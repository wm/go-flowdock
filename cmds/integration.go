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

	// Careful
	// flowsCreate("iora", "wm-test-api", client)
	// flowsUpdate("iora", "wm-test-api", client)

	flowsGet("iora", "culinary-extra", client)
	flowsGetById("iora:culinary-extra", client)
	flowsList(client)

	message := messagesCreate(client)
	messagesComment(client, *message.ID)
	messageList(client)
	// inboxMessage(client)
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
	displayMessageData(*m)

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

// TODO: needs to be fixed (load token from file)
func inboxMessage(client *flowdock.Client) *flowdock.Message {
	opt := &flowdock.InboxCreateOptions{
		Source:            "go-flowdock",
		FromName:          "TeamCity CI",
		Subject:           "IoraHealth/bouncah build #87 has failed!",
		FromAddress:       "build+ok@flowdock.com",
		Link:              "http://wil.io",
		Content:           `
<ul>
	<li>
		<code><a href="https://github.com/IoraHealth/bouncah">IoraHealth/bouncah</a> </code> build #100 has passed!
	</li>
	<li>
		Branch: <code>production</code>
	</li>
	<li>
		Latest commit: <code><a href=\"https://github.com/IoraHealth/bouncah/commit/b35ceee756b579af5e633e8af18b513f7d39470f\">b35ceee</a></code> by <a href=\"mailto:wmernagh@gmail.com\">Will Mernagh</a>
	</li>
	<li>
		Change view: https://github.com/IoraHealth/bouncah/compare/b8036ddd8cd2...b35ceee756b5
	</li>
	<li>
		Build details: https://magnum.travis-ci.com/IoraHealth/bouncah/builds/2104125
	</li>
</ul>
		`,
		Tags:         []string{"fail", "CI", "87"},
	}
	m, _, err := client.Inbox.Create("SOME_TOKEN", opt)
	if err != nil {
		log.Fatal("Get:", err)
	}

	return m
}
