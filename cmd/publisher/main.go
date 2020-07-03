package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dmitsh/pubsub/pkg/pubsub"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var addr, msg string
	a := kingpin.New(filepath.Base(os.Args[0]), "Publisher")
	a.HelpFlag.Short('h')
	a.Flag("address", "Server address.").Short('a').Default("localhost:8080").StringVar(&addr)
	a.Flag("message", "Message to publish.").Short('m').StringVar(&msg)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing commandline arguments: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("Publishing %q to %s", msg, addr)
	if err := pubsub.Publish(addr, msg); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
