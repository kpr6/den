package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
)

var (
	usage = `Usage: den <operation> [options]
Operations:
	get: get the key value from secret den
	set: store key value pair into secret den
	del: delete key value pair from den
	list: list all kv pairs

Options:
	-k	key name
	-v	value name
`
	getUsage = `Usage: den get -k <key-name>`
	setUsage = `Usage: den set -k <key-name> -v <secret-value>`
	delUsage = `Usage: den del -k <key-name>`
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)

	getCmd.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(getUsage))
	}
	setCmd.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(setUsage))
	}
	delCmd.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(delUsage))
	}
	// exit if subcommands aren't passed
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.Parse()
	// get home dir of the user to store the hidden secrets backend file
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	v := config("youjustcantguessit", usr.HomeDir+"/.secrets")

	switch os.Args[1] {
	case "set":
		key := setCmd.String("k", "", "")
		val := setCmd.String("v", "", "")
		setCmd.Parse(os.Args[2:])
		err := v.set(*key, *val)
		if err != nil {
			panic(err)
		}
	case "get":
		key := getCmd.String("k", "", "")
		getCmd.Parse(os.Args[2:])
		value, err := v.get(*key)
		if err != nil {
			panic(err)
		}
		fmt.Println(value)
	case "del":
		key := delCmd.String("k", "", "")
		delCmd.Parse(os.Args[2:])
		err := v.del(*key)
		if err != nil {
			panic(err)
		}
	case "list":
		err := v.list()
		if err != nil {
			panic(err)
		}
	}
}
