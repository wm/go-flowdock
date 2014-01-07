// A command to test against the real service
package main

import (
	"fmt"
	"github.com/wm/go-flowdock/auth"
	"github.com/wm/go-flowdock/flowdock"
	"log"
)

func main() {
	client := flowdock.NewClient(auth.AuthenticationRequest())

	// flowsCreate("iora", "wm-test-api", client)
	// flowsUpdate("iora", "wm-test-api", client)
	flowsGet("iora", "culinary-extra", client)
	flowsGetById("iora:culinary-extra", client)
	flowsList(client)

	message := messagesCreate(client)
	messagesComment(client, *message.ID)
}

func flowsCreate(org, name string, client *flowdock.Client) {
	opt := &flowdock.FlowsCreateOptions{Name: name}
	_, _, err := client.Flows.Create(org, opt)
	if err != nil {
		log.Fatal("Get:", err)
	}
	flowsGet("iora", name, client)
}

func flowsUpdate(org, name string, client *flowdock.Client) {
	disable := true
	flow := &flowdock.Flow{Disabled: &disable}
	flow, _, err := client.Flows.Update(org, name, flow)
	displayFlowData(*flow)
	if err != nil {
		log.Fatal("Get:", err)
	}
}

func flowsGet(org, name string, client *flowdock.Client) {
	flow, _, err := client.Flows.Get(org, name)
	if err != nil {
		log.Fatal("Get:", err)
	}
	displayFlowData(*flow)
}

func flowsGetById(id string, client *flowdock.Client) {
	flow, _, err := client.Flows.GetById(id)

	if err != nil {
		log.Fatal("Get:", err)
	}
	displayFlowData(*flow)
}

func flowsList(client *flowdock.Client) {
	opt := flowdock.FlowsListOptions{User: true}
	flows, _, err := client.Flows.List(true, &opt)

	if err != nil {
		log.Fatal("Get:", err)
	}

	for _, flow := range flows {
		displayFlowData(flow)
	}
}

func displayFlowData(flow flowdock.Flow) {
	org := flow.Organization
	fmt.Println("Flow:", *flow.Id, *flow.Name, *org.Name, *flow.Url)
}

func messagesCreate(client *flowdock.Client) *flowdock.Message {
	opt := &flowdock.MessagesCreateOptions{FlowID: "iora:egg",
		Event: "message",
		Content: "Howdy-Doo @dd #awesome",
		Tags:  []string{"test", ":#api:", "@wm"},
	}
	m, _, err := client.Messages.Create(opt)
	if err != nil {
		log.Fatal("Get:", err)
	}
	fmt.Println("Message", m)

	return m
}

func messagesComment(client *flowdock.Client, messageID int) {
	opt := &flowdock.MessagesCreateOptions{FlowID: "iora:egg",
	    MessageID: messageID,
		Event: "comment",
		Content: "Commenting yo!",
	}
	m, _, err := client.Messages.CreateComment(opt)
	if err != nil {
		log.Fatal("Get:", err)
	}
	fmt.Println("Message", m)
}
