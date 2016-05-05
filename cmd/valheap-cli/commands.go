package main

import (
	"fmt"
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
