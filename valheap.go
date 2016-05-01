package main

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

type DB struct {
	*bolt.DB
}

var (
	userBucket      = []byte(`users`)
	valueBucket     = []byte(`values`)
	ErrUnauthorized = errors.New("Unauthorized")
)

func (db DB) Get(key string) (val []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(valueBucket)
		val = bucket.Get([]byte(key))
		return nil
	})
	return
}

func (db DB) Put(key string, val []byte) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(valueBucket)
		return bucket.Put([]byte(key), val)
	})
	return
}

func EnsureBuckets(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(valueBucket)
		if err != nil {
			return err
		}
		users, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return err
		}
		rootData := users.Get([]byte(`root`))
		if rootData == nil {
			log.Println("Setting up root user (password 'toor', replace it immediately)")
			err = users.Put([]byte(`root`), DefaultRoot.Marshal())
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	EnsureBuckets(db)
	log.Fatal(http.ListenAndServe(":8080", DB{db}.ServeMux()))
}
