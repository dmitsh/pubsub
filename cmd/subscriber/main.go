package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/dmitsh/pubsub/pkg/pubsub"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var n int
	var addr string
	a := kingpin.New(filepath.Base(os.Args[0]), "Subscriber")
	a.HelpFlag.Short('h')
	a.Flag("count", "Number of subscribers.").Short('n').Default("1").IntVar(&n)
	a.Flag("address", "Server address.").Short('a').Default("localhost:8080").StringVar(&addr)
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing commandline arguments: %s", err.Error())
		os.Exit(1)
	}

	data := make(chan string)

	for i := 0; i < n; i++ {
		go func() {
			if err := pubsub.Subscriber(addr, data); err != nil {
				log.Println(err.Error())
			}
		}()
	}

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-term:
			log.Printf("Exit")
			return
		case msg := <-data:
			log.Printf("Received %q", msg)
		}
	}
}
