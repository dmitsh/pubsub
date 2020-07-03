package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dmitsh/pubsub/pkg/pubsub"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var port int
	a := kingpin.New(filepath.Base(os.Args[0]), "PubSub server")
	a.HelpFlag.Short('h')
	a.Flag("port", "Listening port.").Short('p').Default("8080").IntVar(&port)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing commandline arguments: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("Starting server on port %d", port)
	pubsub.NewServer(port).Start()
}
