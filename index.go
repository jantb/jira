package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/boltdb/bolt"
	"github.com/bradfitz/slice"
)

//SearchIndex state
type SearchIndex struct {
	db *bolt.DB
}

// Open a search index
func Open() SearchIndex {
	datastore := SearchIndex{}
	datastore.db = getDb()
	return datastore
}

type similaritystruct struct {
	Key        string
	Similarity float64
}

// Index add issue to index
func (d *SearchIndex) Index(key string, data interface{}) error {
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

// IndexConfluence add confluence page to index
func (d *SearchIndex) IndexConfluence(key string, data interface{}) error {
	datab, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("confluence"))
		b.Put([]byte(key), datab)
		return nil
	})

	return err
}

// IndexSearch add search to index
func (d *SearchIndex) IndexSearch(data interface{}) ([]similaritystruct, error) {
	datab, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("search"))
		b.Put([]byte("search"), datab)
		return nil
	})
	d.calculateTfIdf()
	d.calculateSimularities("search")
	similarities, _ := d.getSimularities("search")
	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("search"))
		b.Delete([]byte("search"))
		return nil
	})
	return similarities, err
}

// Res from the index
type Res struct {
	key   string
	value string
}

// Clear the index
func (d *SearchIndex) Clear() {
	d.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("store"))
		tx.DeleteBucket([]byte("confluence"))
		tx.DeleteBucket([]byte("search"))
		tx.DeleteBucket([]byte("similarDocuments"))
		tx.DeleteBucket([]byte("tfcache"))
		return nil
	})
	return
}

// SearchAllMatching search in the index
func (d *SearchIndex) SearchAllMatching(count int) ([]Res, error) {
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

// SearchAllMatchingSubString searh for substring in the index
func (d *SearchIndex) SearchAllMatchingSubString(s string) ([]Res, error) {
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

func (d *SearchIndex) getKey(key string) (Res, error) {
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

func (d *SearchIndex) getConfluenceKey(key string) (Res, error) {
	var res Res
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("confluence"))
		get := b.Get([]byte(key))
		res.key = key
		res.value = string(get)
		return nil
	})
	return res, err
}

func (d *SearchIndex) getSimularities(key string) ([]similaritystruct, error) {
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

func getStringFromIssue(issue jira.Issue) string {

	comments := ""
	if issue.Fields.Comments != nil {
		for _, comment := range issue.Fields.Comments.Comments {
			if len(comments) != 0 {
				comments += " "
			}
			comments += comment.Body
		}
	}
	return fmt.Sprintf("%s %s %s", issue.Fields.Summary, issue.Fields.Description, comments)
}
func (d *SearchIndex) calculateTfIdf() error {
	if tfidfcache == nil {
		m := make(map[string]string)
		err := d.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("store"))
			b.ForEach(func(k, v []byte) error {
				issue := jira.Issue{}
				json.Unmarshal([]byte(v), &issue)
				m[string(k)] = getStringFromIssue(issue)
				return nil
			})
			return nil
		})

		err = d.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("confluence"))
			b.ForEach(func(k, v []byte) error {
				page := Page{}
				json.Unmarshal([]byte(v), &page)
				m[string(k)] = page.Body
				return nil
			})
			return nil
		})

		err = d.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("search"))
			b.ForEach(func(k, v []byte) error {
				page := Page{}
				json.Unmarshal([]byte(v), &page)
				m[string(k)] = string(v)
				return nil
			})
			return nil
		})

		if err != nil {
			return err
		}

		tfidfcache = tfidfMap(m)
		bytes, _ := json.Marshal(tfidfcache)
		err = d.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("tfcache"))
			b.Put([]byte("tf"), bytes)
			return nil
		})
	}
	return nil
}

func (d *SearchIndex) calculateSimularities(key string) error {
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tfcache"))
		bytes := b.Get([]byte("tf"))
		if bytes == nil {
			d.calculateTfIdf()
		}
		json.Unmarshal(bytes, &tfidfcache)
		return nil
	})

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

	similaritiesb, err := json.Marshal(similarities[:Min(200, len(similarities))])
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
		tx.CreateBucketIfNotExists([]byte("confluence"))
		tx.CreateBucketIfNotExists([]byte("search"))
		tx.CreateBucketIfNotExists([]byte("similarDocuments"))
		tx.CreateBucketIfNotExists([]byte("tfcache"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return dbs
}
