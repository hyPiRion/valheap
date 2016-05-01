package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
)

const banner = `%s is a tool to talk to a valheap server easily

The following commands are available:

init    - (re)Set up your configuration
chgwpwd - Change your valheap password
get     - Get a key from valheap and print to stdout
put     - Put/update a key to valheap from stdin
adduser - Adds a user to valheap (must be root)
rmuser  - Removes a user from valheap (must be root)

The environment variable VALHEAP_CLI_FILE can be set to override the
default valheap file location, which is $HOME/.valheap-cli.json.
`

var knownCommand map[string]bool

func init() {
	knownCommand = map[string]bool{
		"init":    true,
		"chgpwd":  true,
		"get":     true,
		"put":     true,
		"adduser": true,
		"rmuser":  true,
	}
}

type Config struct {
	Server   string
	Username string
	Password []byte
}

var cfg Config

func configPath() string {
	path := os.Getenv("VALHEAP_CLI_FILE")
	if path == "" {
		var err error
		path, err = homedir.Expand("~/.valheap-cli.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to detect home directory: %s\n", err)
			os.Exit(1)
		}
	}
	return path
}

func readInit() error {
	if _, err := os.Stat(configPath()); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Valheap configuration is not set up")
		fmt.Fprintf(os.Stderr, "Do `%s init` to setup valheap-cli\n", os.Args[0])
		os.Exit(1)
	}
	return tryReadInit()
}

func tryReadInit() error {
	bs, err := ioutil.ReadFile(configPath())
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &cfg)
}

func makeInit() {
	tryReadInit()

	scanner := bufio.NewScanner(os.Stdin)
	if cfg.Server != "" {
		fmt.Printf("  Old host location: %s\n", cfg.Server)
	}
	fmt.Print("Enter host location: ")
	scanner.Scan()
	server := scanner.Text()
	if cfg.Username != "" {
		fmt.Printf("  Old username: %s\n", cfg.Username)
	}
	fmt.Print("Enter username: ")
	scanner.Scan()
	uname := scanner.Text()

	fmt.Print("Enter password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Error reading password: %s\n", err)
		os.Exit(1)
	}

	path := configPath()
	err = os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	newConf := Config{
		Server:   server,
		Username: uname,
		Password: pass,
	}

	bs, _ := json.Marshal(newConf)

	err = ioutil.WriteFile(path, bs, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// try to connect later on
}

func main() {
	if len(os.Args) == 1 || !knownCommand[os.Args[1]] {
		fmt.Printf(banner, os.Args[0])
		os.Exit(1)
	}
	if os.Args[1] == "init" {
		makeInit()
	}
	err := readInit()
	if err != nil {
		panic(err)
	}
}
