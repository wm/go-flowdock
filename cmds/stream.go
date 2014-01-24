package main

import (
	"fmt"
	"github.com/wm/go-flowdock/auth"
	"github.com/wm/go-flowdock/flowdock"
	"log"
	"code.google.com/p/goauth2/oauth"
)

func main() {
	httpClient := auth.AuthenticationRequest()
	token, _ := oauth.CacheFile("cache.json").Token()

	client := flowdock.NewClient(httpClient)

	messageList(client)
	messageStream(client,token.AccessToken)

	fmt.Println("Waiting for event")
}

func messageStream(client *flowdock.Client, token string) {
	stream, es, _ := client.Messages.Stream(token, "iora", "tech-stuff")
	stream1, es1, _ := client.Messages.Stream(token, "iora", "technical-discussions")
	defer es.Close()
	defer es1.Close()

	for {
		select {
		case msg := <-stream:
			displayMessageData(msg, "wc")
		case msg1 := <-stream1:
			displayMessageData(msg1, "td")
		}
	}
}

func messageList(client *flowdock.Client) {
	opt := flowdock.MessagesListOptions{Limit: 100}
	messages, _, err := client.Messages.List("iora", "tech-stuff", &opt)

	if err != nil {
		log.Fatal("Get:", err)
	}

	for _, msg := range messages {
		displayMessageData(msg, "wc")
	}
}

func displayMessageData(msg flowdock.Message, room string) {
	events := []string{"user-edit", "file", "activity.user", "mail", "zendesk", "twitter", "tag-change"}
	if stringNotInSlice(*msg.Event, events) {
		fmt.Println("\nMSG:", room, *msg.ID, *msg.Event, msg.Content())
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func stringNotInSlice(a string, list []string) bool {
	return !stringInSlice(a, list)
}
