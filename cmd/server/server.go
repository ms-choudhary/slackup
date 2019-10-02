package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ms-choudhary/slackup/pkg/server"
	"github.com/ms-choudhary/slackup/pkg/store"
)

var (
	port    = flag.Uint("port", 8080, "The port to listen on.  Default 8080.")
	address = flag.String("address", "127.0.0.1", "The address on the local server to listen to. Default 127.0.0.1")
	dbFile  = flag.String("dbfile", "", "Sqlite database file")
)

func main() {
	flag.Parse()

	if *dbFile == "" {
		fmt.Fprintln(os.Stderr, "missing required option: dbfile")
		flag.Usage()
		os.Exit(1)
	}

	store, err := store.Init(*dbFile)
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}
	defer store.Close()

	server := &server.Server{
		Store: store,
	}

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *address, *port),
		Handler: server,
	}

	log.Printf("starting up server at %s:%d ...", *address, *port)

	log.Fatal(s.ListenAndServe())
}
