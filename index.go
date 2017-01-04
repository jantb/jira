package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/bradfitz/slice"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type searchIndex struct {
	db *bolt.DB
}

func Open() searchIndex {
	datastore := searchIndex{}
	datastore.db = getDb()
	//v, _ := json.MarshalIndent(datastore.db.Stats(), "", "    ")
	//fmt.Println(string(v))
	return datastore
}

type similaritystruct struct {
	Key        string
	Similarity float64
}

func (d *searchIndex) Index(key string, data interface{}) error {
	datab, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		b.Put([]byte(key), datab)
		return nil
	})

	return err
}

type Res struct {
	key   string
	value string
}

func (d *searchIndex) Clear() {
	d.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("store"))
		tx.DeleteBucket([]byte("similarDocuments"))
		tx.DeleteBucket([]byte("tfcache"))
		return nil
	})
	return
}

func (d *searchIndex) SearchAllMatching(count int) ([]Res, error) {
	var res []Res
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		c := b.Cursor()

		for k, v := c.First(); k != nil && count > 0; k, v = c.Next() {
			res = append(res, Res{key: string(k), value: string(v)})
			count--
		}
		return nil
	})
	return res, err
}

func (d *searchIndex) SearchAllMatchingSubString(s string) ([]Res, error) {
	var res []Res
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Index(string(k), s) > -1 {
				res = append(res, Res{key: string(k), value: string(v)})
			}
		}
		return nil
	})
	return res, err
}

func (d *searchIndex) getKey(key string) (Res, error) {
	var res Res
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		get := b.Get([]byte(key))
		res.key = key
		res.value = string(get)
		return nil
	})
	return res, err
}

func (d *searchIndex) getSimularities(key string) ([]similaritystruct, error) {
	similarities := make([]similaritystruct, 0)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("similarDocuments"))
		by := b.Get([]byte(key))
		json.Unmarshal(by, &similarities)
		return nil
	})
	if err != nil {
		return similarities, err
	}

	return similarities, nil
}

var tfidfcache map[string]map[string]float64

func (d *searchIndex) calculateTfIdf() (err error) {
	if tfidfcache == nil {
		t := time.Now()
		m := make(map[string]string)
		err := d.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("store"))
			b.ForEach(func(k, v []byte) error {
				m[string(k)] = string(v)
				return nil
			})
			return nil
		})
		if err != nil {
			return err
		}

		fmt.Print("generating tf-idf map... ")
		tfidfcache = tfidfMap(m)
		bytes, _ := json.Marshal(tfidfcache)
		err = d.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("tfcache"))
			b.Put([]byte("tf"), bytes)
			return nil
		})
		fmt.Println(time.Now().Sub(t))
	}
	return nil
}

func (d *searchIndex) calculateSimularities(key, data string) error {
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tfcache"))
		bytes := b.Get([]byte("tf"))
		if bytes == nil {
			d.calculateTfIdf()
		}
		json.Unmarshal(bytes, &tfidfcache)
		return nil
	})

	t := time.Now()
	fmt.Print("                                           \rgenerating similarities map... " + key + " ")
	tfidfdata := tfidfcache[key]
	var similarities []similaritystruct
	for k, value := range tfidfcache {
		if k == key {
			continue
		}
		similarities = append(similarities, similaritystruct{
			Key:        k,
			Similarity: similarity(tfidfdata, value),
		})
	}
	slice.Sort(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	similaritiesb, err := json.Marshal(similarities[:20])
	if err != nil {
		return err
	}
	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("similarDocuments"))
		b.Put([]byte(key), similaritiesb)
		return nil
	})
	fmt.Print(time.Now().Sub(t))
	return nil
}

func getDb() *bolt.DB {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := bolt.Open(filepath.Join(usr.HomeDir, ".jira.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = dbs.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("store"))
		tx.CreateBucketIfNotExists([]byte("similarDocuments"))
		tx.CreateBucketIfNotExists([]byte("tfcache"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return dbs
}
