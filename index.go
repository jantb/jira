package main

import (
	"github.com/boltdb/bolt"
	"path/filepath"
	"log"
	"os/user"
	"encoding/json"
	"github.com/bradfitz/slice"
	"fmt"
	"time"
)

type searchIndex struct {
	db *bolt.DB
}

func Open() searchIndex {
	datastore := searchIndex{}
	datastore.db = getDb()
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
func (d *searchIndex) SearchAllMatching(count int) ([][]byte, error) {
	var res [][]byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		b.ForEach(func(k, v []byte) error {
			res = append(res, v)
			count--
			if count == 0 {
				return nil
			}
			return nil
		})
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

func (d *searchIndex) calculateSimularities(key, data string) error {
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

		m[key] = data
		fmt.Println("generating tf-idf map...")
		tfidfcache = tfidfMap(m)
		fmt.Println()
		fmt.Println(time.Now().Sub(t))
	}

	tfidfdata := tfidfcache[key]
	var similarities []similaritystruct
	for k, value := range tfidfcache {
		if k == key {
			continue
		}
		similarities = append(similarities, similaritystruct{
			Key:k,
			Similarity:similarity(tfidfdata, value),
		})
	}
	slice.Sort(similarities, func(i, j int) bool {
		return similarities[i].Similarity < similarities[j].Similarity
	})

	similaritiesb, err := json.Marshal(similarities)
	if err != nil {
		return err
	}
	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("similarDocuments"))
		b.Put([]byte(key), similaritiesb)
		return nil
	})

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
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return dbs
}
