package main

import (
	"github.com/boltdb/bolt"
	"path/filepath"
	"log"
	"os/user"
	"encoding/json"
	"github.com/bradfitz/slice"
	"fmt"
)

type datastore struct {
	db *bolt.DB
}

func Open() datastore {
	datastore := datastore{}
	datastore.db = getDb()
	return datastore
}

type similaritystruct struct {
	Key        string
	Similarity float64
}

func (d *datastore) Index(key string, data interface{}) error {
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
func (d *datastore) SearchAllMatching(count int) ([][]byte,error) {
	var res [][]byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		b.ForEach(func(k, v []byte) error {
			res = append(res,v)
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

func (d *datastore) getSimularities(key string) ([]similaritystruct, error) {
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
func (d *datastore) calculateSimularities(key, data string) error {
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

	tfidf := tfidfMap(m)
	tfidfdata := tfidf[key]

	var similarities []similaritystruct
	for k, value := range tfidf {
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
