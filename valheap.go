package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

func main() {
	var dbpath, certFile, keyFile string
	var help bool
	var port int
	flag.StringVar(&dbpath, "db", "valheap.db", "Path to the bolt DB file to use")
	flag.IntVar(&port, "port", 8080, "The port to listen on HTTP requests")
	flag.BoolVar(&help, "help", false, "Prints this help message")
	flag.StringVar(&certFile, "cert", "", "The path to the TLS certificate to use")
	flag.StringVar(&keyFile, "key", "", "The path to the TLS private key to use")
	flag.Parse()

	if help {
		fmt.Println(`Valheap is an HTTP key/value storage with basic auth

Usage: ./valheap [options]

Where options may be:`)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if (certFile == "" || keyFile == "") && keyFile != certFile {
		fmt.Fprintln(os.Stderr, "Both -cert and -key must be specified to use TLS")
		os.Exit(1)
	}

	log.Infof("Opening database file %s", dbpath)
	db, err := bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	EnsureBuckets(db)

	addr := fmt.Sprintf(":%d", port)

	log.Infof("Now listening on port %d", port)
	if certFile != "" {
		err = http.ListenAndServeTLS(addr, certFile, keyFile, DB{db}.ServeMux())
	} else {
		log.Warning("Not using TLS. If you want to be secure, either enable it or put this behind nginx or something similar")
		err = http.ListenAndServe(addr, DB{db}.ServeMux())
	}
	log.Fatal(err)
}
