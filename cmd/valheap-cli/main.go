package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
)

const banner = `%s is a tool to talk to a valheap server easily

The following commands are available:

init       (re)Set up your configuration
chgwpwd    Change your valheap password
get        Get a key from valheap and print to stdout
put        Put/update a key to valheap from stdin
delete     Deletes a key from valheap
list       Lists all keys in valheap with the provided prefix
adduser    Adds a user to valheap (must be root)
rmuser     Removes a user from valheap (root only)
listusers  Lists all users in valheap (root only)
backup     Backups the database to the provided file (root only)

The environment variable VALHEAP_CLI_FILE can be set to override the
default valheap file location, which is $HOME/.valheap-cli.json.
`

var knownCommand map[string]func(string)

func init() {
	knownCommand = map[string]func(string){
		"init":      Get, // yeah yeah
		"get":       Get,
		"put":       Put,
		"delete":    Delete,
		"adduser":   AddUser,
		"rmuser":    RmUser,
		"chgpwd":    Get, // dummy
		"list":      List,
		"listusers": List,
		"backup":    Backup,
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
		fmt.Printf("  Old server URL: %s (just press enter to keep it)\n", cfg.Server)
	}
	fmt.Print("Enter server URL: ")
	scanner.Scan()
	server := scanner.Text()
	if server == "" {
		server = cfg.Server
	}
	_, err := url.Parse(server)
	if err != nil {
		fmt.Printf("Bad URL: %s\n", err)
		os.Exit(1)
	}
	if cfg.Username != "" {
		fmt.Printf("  Old username: %s (just press enter to keep it)\n", cfg.Username)
	}
	fmt.Print("Enter username: ")
	scanner.Scan()
	uname := scanner.Text()
	if uname == "" {
		uname = cfg.Username
	}

	fmt.Print("Enter password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Error reading password: %s\n", err)
		os.Exit(1)
	}

	// verify that information is correct before storing it.
	req, err := http.NewRequest("GET", server, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(uname, string(pass))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Unable to verify correct info: %s\n", err)
		os.Exit(1)
	}
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		fmt.Println("Incorrect username/password, please try again")
		os.Exit(1)
	case http.StatusNotFound:
	default:
		fmt.Printf("Unexpected error code from server: %d\n", resp.StatusCode)
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
	os.Exit(0)
}

func main() {
	// Well, this is an interesting piece of spaghetti
	if len(os.Args) == 1 || knownCommand[os.Args[1]] == nil {
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
	if os.Args[1] == "chgpwd" {
		ChgPwd()
		os.Exit(0)
	}
	if os.Args[1] == "listusers" {
		ListUsers()
		os.Exit(0)
	}
	if os.Args[1] == "list" {
		if len(os.Args) > 3 {
			fmt.Fprintf(os.Stderr, "%s expects 0 or 1 argument in\n", os.Args[1])
			os.Exit(1)
		}
		prefix := ""
		if len(os.Args) == 3 {
			prefix = os.Args[2]
		}
		List(prefix)
		os.Exit(0)
	}
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "%s expects exactly 1 argument in\n", os.Args[1])
		os.Exit(1)
	}
	knownCommand[os.Args[1]](os.Args[2])
}
