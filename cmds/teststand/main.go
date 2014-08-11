package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var usage = `
Usage:
    teststand serve
    teststand client [--name=me]
    teststand take [--name=me]              Take the mutex. Prints status while waiting.
    teststand give [--name=me]              Give the mutex to someone waiting.
    teststand message [--name=me] "Message" Message others waiting on the mutex know.
    teststand ds [--name=me] enable         Set the ds state to either enable, disable or resync.
`

// main handles the arguments and runs the proper function based off the args.
func main() {
	if len(os.Args) <= 1 {
		fmt.Println(usage)
		return
	}

	name := getArg("name", os.Getenv("USER"))
	addr := getArg("addr", "10.1.90.2:8080")

	switch os.Args[1] {
	case "take", "--take", "t":
		err := take(addr, name)
		if err != nil {
			log.Fatal(err)
		}

	case "give", "--give", "g":
		err := give(addr, name)
		if err != nil {
			log.Fatal(err)
		}

	case "message", "--message", "m":
		if len(os.Args) < 3 {
			log.Fatal("Need to include a message")
		}
		err := message(addr, name, strings.Join(os.Args[2:], " "))
		if err != nil {
			log.Fatal(err)
		}

	case "ds", "--ds", "d":
		if len(os.Args) != 3 {
			log.Fatal("Need to specify a state")
		}
		err := dsControl(addr, name, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

	case "--serve", "--server", "serve", "server":
		addr := ":8080"
		if len(os.Args) == 3 {
			addr = os.Args[2]
		}
		serve(addr)

	case "--help", "-help", "-h", "h", "help", "usage":
		// Explicitly intended help flags
		fallthrough
	default:
		fmt.Println(usage)
	}
}

// getArg returns the value of the argument `name` or the default value.
func getArg(name, def string) string {
	var value string
	var index = -1
	for i, s := range os.Args {
		if strings.HasPrefix(s, "--"+name+"=") {
			index = i
		}
	}
	if index == -1 {
		return def
	}
	value = strings.TrimPrefix(os.Args[index], "--"+name+"=")
	os.Args = append(os.Args[:index], os.Args[1+index:]...)
	return value
}
