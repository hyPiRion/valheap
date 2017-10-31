package main

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	*bolt.DB
}

var (
	userBucket      = []byte(`users`)
	valueBucket     = []byte(`values`)
	ErrUnauthorized = errors.New("Unauthorized")
)

func copyBytes(bs []byte) []byte {
	if bs == nil {
		return nil
	}
	ret := make([]byte, len(bs), len(bs))
	copy(ret, bs)
	return ret
}

func (db DB) Get(key string) (val []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(valueBucket)
		val = copyBytes(bucket.Get([]byte(key)))
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

func (db DB) Delete(key string) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(valueBucket)
		return bucket.Delete([]byte(key))
	})
	return
}

func (db DB) List(prefixString string) (keys [][]byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(valueBucket).Cursor()
		prefix := []byte(prefixString)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			keys = append(keys, copyBytes(k))
		}
		return nil
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
