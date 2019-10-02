package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ms-choudhary/slackup/pkg/api"
	"github.com/ms-choudhary/slackup/pkg/store"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("invalid arguments: %v", os.Args)
	}
	dbFile := os.Args[1]
	db, err := store.Init(dbFile)
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}
	defer db.Close()

	id, err := store.GetChannel(db, "scripbox", "ops-incident")

	if err != nil {
		log.Fatalf("failed to get channel: %v", err)
	}

	msgs := []api.Message{
		api.Message{
			User: "mohit",
			Text: "hello world",
			Ts:   "123",
			Comments: []api.Message{
				api.Message{User: "parthesh", Text: "howdy", Ts: "124"},
			},
		},
	}

	err = store.Insert(db, id, msgs)
	if err != nil {
		log.Fatalf("failed to insert msgs: %v", err)
	}
	msgs, err = store.Query(db, id, store.Filter{})
	if err != nil {
		log.Fatalf("failed to query msgs: %v", err)
	}
	fmt.Println(msgs)
}
