package main

import (
	"github.com/boltdb/bolt"
	"path/filepath"
	"log"
	"os/user"
)

type datastore struct{
	db *bolt.DB
}

func Open() datastore{
	datastore := datastore{}
	datastore.db = getDb()
	return datastore
}
func (d *datastore) index(data string) error{
	var documents []string
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("store"))
		b.ForEach(func(k,v []byte) {
			documents = append(documents, string(v))
		})
		return nil
	})



	return err
}

func getDb()*bolt.DB{
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
		tx.CreateBucketIfNotExists([]byte("tfidf"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return dbs
}
