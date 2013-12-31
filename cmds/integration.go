// A command to test against the real service
package main

import (
	"fmt"
	"log"
	"github.com/wm/go-flowdock/flowdock"
	"github.com/wm/go-flowdock/auth"
)

func main() {
	client := flowdock.NewClient(auth.AuthenticationRequest())

	// flowsCreate("iora", "wm-test-api", client)
	// flowsUpdate("iora", "wm-test-api", client)
	flowsGet("iora", "culinary-extra", client)
	flowsGetById("iora:culinary-extra", client)
	flowsList(client)
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
	flow := &flowdock.Flow{Disabled: true}
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
	fmt.Println("Flow:", flow.Id, flow.Name, org.Name, flow.Url)
}

