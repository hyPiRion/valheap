# Valheap

Valheap is an HTTP key/value store to store small things, with basic
authentication on top. It's designed to be easy to use, easy to maintain/set up
and easy to understand.

It's **not** intended to be

* Scalable – The server application is one one machine, and one machine only.
* Extremely secure – If someone is able to get into the machine valheap runs on,
  chances are everything in the store has been compromised. (But perhaps you
  have bigger issues if someone is able to get into it)
* A database – You can only get, put, delete and list values. You cannot iterate
  over them, do queries on them and so on, without doing it manually and
  non-transactionally.
* Efficient – You should probably not hook this up to systems with "millions of
  requests per second" (but you could try if you want to).
* Used to store big things. The server should be able to keep the thing you
  store in memory.

## Quickstart

Do `go get github.com/hyPiRion/valheap/...` to install both valheap and
valheap-cli.

To run a local version of valheap, start valheap in a terminal window like so:

```
valheap
```


It should emit something like this:

```
INFO[0000] Opening database file valheap.db
INFO[0000] Setting up root user (password 'toor', replace it immediately)
INFO[0000] Now listening on port 8080
WARN[0000] Not using TLS. If you want to be secure, either enable it or put this
  behind nginx or something similar
```

To communicate with it, you can use `valheap-cli` in another terminal window.
You can issue get, put and delete, but first we have to configure valheap-cli:

```
$ valheap-cli get 100
Valheap configuration is not set up
Do `valheap-cli init` to setup valheap-cli
$ valheap-cli init
Enter server URL: http://localhost:8080
Enter username: root
Enter password: toor # not shown
```

We should probably start off by changing the root password though, as instructed
by valheap itself. This can be done through `valheap-cli chgpwd`:

```
$ valheap-cli chgpwd
Enter new password: 
$ ..
```

You don't have to specify the password whenever you call valheap-cli: It is
stored in the file `~/.valheap-cli.json` (in plaintext). This also means that
people who temporarily use your machine can get your password, the same applies
if your machine is stolen and you don't have the hard drive encrypted. (You can
use the raw HTTP endpoint instead if you prefer to not store the password.)

### Adding, Reading and Deleting Keys

`valheap-cli` uses stdin and stdout for gets, puts and deletes: 

```shell
$ res=$(valheap-cli get foo)
Not Found
# Also the exit status is 1
$ echo $res

$ echo 'bar' | valheap-cli put foo
$ valheap-cli get foo
bar
$ res=$(valheap-cli get foo)
$ echo $res
bar
$ valheap-cli delete foo
Key foo deleted
```

There are no limitations to a name of a key, just beware of URL encoding.
`valheap-cli will` url encode keys, but only the things that have to be encoded.
For example, although 

### Listing Keys

By using the `list` command, you can list all the keys in the project, or just a
subset by specifying a prefix. The keys are separated by a newline, but be aware
that keys themselves contain newlines. There is no attempt at handling this on
the client nor the server:

```shell
$ valheap-cli list
bar
foo
$ valheap-cli list f
foo
$ valheap-cli list dev
# Nothing to be printed
$
```

### Adding and Removing Users

Only the root user can add and remove arbitrary users, and root cannot be
removed. Users can remove themselves, but cannot remove anyone else or add
anyone.

To add a user via valheap-cli, use the adduser command with the name of the
user:

```shell
$ valheap-cli adduser fatimah
Enter password for fatimah: 
User updated/added
```

If that's your official account, you can change it by using `valheap-cli init`
again:

```shell
$ valheap-cli init
  Old server URL: http://localhost:8080 (just press enter to keep it)
Enter server URL: 
  Old username: root (just press enter to keep it)
Enter username: fatimah
Enter password:
```

The root user can also remove users by calling `rmuser`, or list them:

```shell
$ valheap-cli rmuser benjamin
User removed
```

```shell
$ valheap-cli listusers
fatimah
root
trevor
```

### Performing Backups

Root users can perform backups via the backup command. A nonexisting file path
must be provided to store the file:

```shell
$ valheap-cli backup mybackup.db
Backup saved in mybackup.db
```

The backup will be stored with the permissions 0600 (only you can read or write
the file)

## Deploying

TODO

## License

Copyright © 2016 Jean Niklas L'orange

Distributed under the BSD 3-clause license, which is available in the file
LICENSE.
