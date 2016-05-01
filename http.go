package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

func (db DB) ServeMux() *http.ServeMux {
	sm := http.NewServeMux()
	sm.HandleFunc("/user/", db.HttpAuth(db.HttpHandleUser))
	sm.HandleFunc("/val/", db.HttpAuth(db.HttpGetPut))
	sm.HandleFunc("/", db.HttpAuth(http.NotFound))
	return sm
}

func (db DB) HttpAuth(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uname, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		err := db.View(func(tx *bolt.Tx) error {
			err := AuthorizeUser(tx, uname, pass)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return ErrUnauthorized
			}
			return nil
		})
		switch err {
		case ErrUnauthorized:
		case nil:
			handler(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func (db DB) HttpHandleUser(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/user/")
	if name == "" {
		http.NotFound(w, r) // I guess?
		return
	}
	uname, _, _ := r.BasicAuth()
	switch r.Method {
	case "PUT":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Unable to read request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		u, err := UnmarshalUser(body)
		if err != nil {
			log.Errorf("PUT /user/%s: %s", name, err)
			http.Error(w, `Request must be in JSON on form {"Password": "mypass"}`, http.StatusBadRequest)
			return
		}
		err = db.PutUser(uname, name, u)
		switch err {
		case ErrForbiddenRoot:
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			log.Errorf("Unexpected error adding user: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		case nil:
			io.WriteString(w, "User updated/added\n")
		}
	case "DELETE":
		err := db.RmUser(uname, name)
		switch err {
		case ErrForbiddenRoot:
			http.Error(w, err.Error(), http.StatusForbidden)
		case ErrUserNotExists:
			http.Error(w, "The user does not exists", http.StatusConflict)
		default:
			log.Errorf("Unexpected error removing user: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		case nil:
			io.WriteString(w, "User removed\n")
		}
	default:
		http.NotFound(w, r)
		return
	}
}

func (db DB) HttpGetPut(w http.ResponseWriter, r *http.Request) {
	keyStr := strings.TrimPrefix(r.URL.Path, "/val/")
	switch r.Method {
	case "PUT":
		// read value
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Unable to read request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = db.Put(keyStr, body)
		if err != nil {
			log.Errorf("Unable to put key: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(body)
		if err != nil {
			log.Errorf("Unable to send body to request: %s", err)
		}
	case "GET":
		val, err := db.Get(keyStr)
		if err != nil {
			log.Errorf("Unable to retrieve key %q: %s", keyStr, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if val == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		_, err = w.Write(val)
		if err != nil {
			log.Errorf("Unable to send body to request: %s", err)
		}
	default:
		http.NotFound(w, r)
	}
}
