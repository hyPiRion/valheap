package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func Get(val string) {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/val/%s", u.Path, val)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(cfg.Username, string(cfg.Password))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		os.Stderr.Write(body)
		os.Exit(1)
	}
	_, err = os.Stdout.Write(body)
	if err != nil {
		os.Exit(1)
	}
}

func Put(val string) {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/val/%s", u.Path, val)

	req, err := http.NewRequest("PUT", u.String(), os.Stdin)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(cfg.Username, string(cfg.Password))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		os.Stderr.Write(body)
		os.Exit(1)
	}
}

func Delete(val string) {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/val/%s", u.Path, val)

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(cfg.Username, string(cfg.Password))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		os.Stderr.Write(body)
		os.Exit(1)
	}
	_, err = os.Stdout.Write(body)
	if err != nil {
		os.Exit(1)
	}
}

func List(val string) {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/listvals", u.Path)
	q := u.Query()
	q.Set("prefix", val)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(cfg.Username, string(cfg.Password))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		os.Stderr.Write(body)
		os.Exit(1)
	}
	_, err = os.Stdout.Write(body)
	if err != nil {
		os.Exit(1)
	}
}

func backup(path string) int {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path += "/backup"

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(cfg.Username, string(cfg.Password))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if resp.StatusCode != 200 {
		io.Copy(os.Stderr, resp.Body)
		return 1
	}

	// Open a new file, must not already exist
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println("Backup saved in file", path)
	return 0
}

// Wrap backup to ensure defer call happens
func Backup(path string) {
	os.Exit(backup(path))
}
