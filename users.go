package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 6

var ErrDBCorrupted = errors.New("Unexpected error (database corrupted?)")
var ErrForbiddenRoot = errors.New("Forbidden: Must be root (or the user itself) to do this")
var ErrUserNotExists = errors.New("User does not exist")
var ErrCannotDeleteRoot = errors.New("Cannot delete the root user")

type User struct {
	HashPass string
}

func (u *User) Authorize(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashPass), []byte(pass))
}

func (u User) Marshal() []byte {
	bs, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return bs
}

func UnmarshalRawUser(data []byte) (u *User, err error) {
	err = json.Unmarshal(data, &u)
	return
}

func UnmarshalUser(data []byte) (u *User, err error) {
	dataShape := struct {
		Password string
	}{}
	err = json.Unmarshal(data, &dataShape)
	if err != nil {
		return
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(dataShape.Password), bcryptCost)
	if err != nil {
		return
	}
	return &User{HashPass: string(bs)}, nil
}

var DefaultRoot = User{}

func init() {
	bs, _ := bcrypt.GenerateFromPassword([]byte(`toor`), bcryptCost)
	DefaultRoot.HashPass = string(bs)
}

// Puts a user into the system, potentially replacing it with the existing user.
// Only root and the user itself can modify this value
func (db DB) PutUser(name, putUname string, u *User) error {
	if name != "root" && name != putUname {
		return ErrForbiddenRoot
	}
	return db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(userBucket)
		return users.Put([]byte(putUname), u.Marshal())
	})
}

func (db DB) RmUser(name, toDelete string) error {
	if name != "root" && name != toDelete {
		return ErrForbiddenRoot
	}
	if toDelete == "root" {
		return ErrCannotDeleteRoot
	}
	return db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(userBucket)
		uinfo := users.Get([]byte(toDelete))
		if uinfo == nil {
			return ErrUserNotExists
		}
		return users.Delete([]byte(toDelete))
	})
}

func AuthorizeUser(tx *bolt.Tx, name, pass string) error {
	users := tx.Bucket(userBucket)
	udata := users.Get([]byte(name))
	if udata == nil {
		return fmt.Errorf("No user with username %q", name)
	}
	user, err := UnmarshalRawUser(udata)
	if err != nil {
		return ErrDBCorrupted
	}
	return user.Authorize(pass)
}
