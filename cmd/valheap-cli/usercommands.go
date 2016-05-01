package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/howeyc/gopass"
)

type JsonOutput struct {
	Password string
}

func AddUser(username string) {
	if username == cfg.Username {
		ChgPwd()
		return
	}
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/user/%s", u.Path, url.QueryEscape(username))

	fmt.Printf("Enter password for %s: ", username)
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Error reading password: %s\n", err)
		os.Exit(1)
	}

	bs, err := json.Marshal(JsonOutput{string(pass)})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(bs))
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
	os.Stdout.Write(body)
}

func RmUser(username string) {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/user/%s", u.Path, url.QueryEscape(username))

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
	os.Stdout.Write(body)
}

func ChgPwd() {
	u, err := url.Parse(cfg.Server)
	if err != nil {
		panic(err)
	}
	u.Path = fmt.Sprintf("%s/user/%s", u.Path, url.QueryEscape(cfg.Username))

	fmt.Print("Enter new password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Error reading password: %s\n", err)
		os.Exit(1)
	}

	bs, err := json.Marshal(JsonOutput{string(pass)})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(bs))
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

	path := configPath()
	err = os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg.Password = pass

	bs, _ = json.Marshal(cfg)

	err = ioutil.WriteFile(path, bs, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
	os.Stdout.Write(body)
}
