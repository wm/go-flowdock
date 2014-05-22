package main

import (
	"fmt"
	"github.com/wm/go-flowdock/auth"
	"github.com/wm/go-flowdock/flowdock"
	"log"
)

// Example of searching a flowdock flow for a string and tags
func main() {
	client := flowdock.NewClient(auth.AuthenticationRequest())

	search := "production to production"
	tags := []string{"deployment", "deploy_end", "production", "icis"}
	event := "mail"
	messageSearch(&tags, &event, &search, client)
}

func messageSearch(tags *[]string, event, search *string, client *flowdock.Client) {
	opt := flowdock.MessagesListOptions{Limit: 100, TagMode: "and"}
	if tags != nil {
		opt.Tags = *tags
	}
	if search != nil {
		opt.Search = *search
	}
	if event != nil {
		opt.Event = *event
	}
	messages, _, err := client.Messages.List("iora", "tech-stuff", &opt)

	if err != nil {
		log.Fatal("Get:", err)
	}

	for _, msg := range messages {
		displayMessageData(msg)
	}
	fmt.Println("Count:", len(messages))
}

func displayMessageData(msg flowdock.Message) {
	fmt.Println("MSG:", *msg.Sent, *msg.ID, *msg.Event, *msg.Tags)
}
