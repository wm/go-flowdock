package main

import (
	"fmt"
	"github.com/wm/go-flowdock/auth"
	"github.com/wm/go-flowdock/flowdock"
	"log"
	"github.com/bernerdschaefer/eventsource"
	"time"
	"code.google.com/p/goauth2/oauth"
)

func main() {
	httpClient := auth.AuthenticationRequest()
	token, _ := oauth.CacheFile("cache.json").Token()
	fmt.Println("Token:", token.AccessToken)
	flow := fmt.Sprintf("flows/iora/egg?access_token=%v", token.AccessToken)

	client := flowdock.NewClient(httpClient)

	req, _ := client.NewStreamRequest("GET", flow, nil)
	es := eventsource.New(req, 3*time.Second)
	messageList(client)

	fmt.Println("Waiting for event")

	for i := 0; i < 10; i++ {
		event, err := es.Read()
		if err == nil {
			fmt.Println(event.ID, event.Type, string(event.Data))
			fmt.Println("And the event is:")
		} else {
			fmt.Println("And the error is:", err, req)
		}
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
