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
	fmt.Println("Token:", token.AccessToken)
	// flow := fmt.Sprintf("flows/iora/egg?access_token=%v", token.AccessToken)

	client := flowdock.NewClient(httpClient)

	messageList(client)
	messageStream(client,token.AccessToken)

	fmt.Println("Waiting for event")
}

func messageStream(client *flowdock.Client, token string) {
	stream, _, _ := client.Messages.Stream(token, "iora", "egg")

	for m := range stream {
		fmt.Println(m)
	}
}

func messageList(client *flowdock.Client) {
	opt := flowdock.MessagesListOptions{Limit: 100, Event: "message, comment"}
	messages, _, err := client.Messages.List("iora", "egg", &opt)

	if err != nil {
		log.Fatal("Get:", err)
	}

	for _, msg := range messages {
		displayMessageData(msg)
	}
}

func displayMessageData(msg flowdock.Message) {
	fmt.Println("MSG:", *msg.ID, *msg.Event, msg.Content())
}
