package flowdock

import (
	"fmt"
	"net/http"
)

// FlowsService handles communication with the flow related methods of the
// Flowdock API.
//
// Flowdock API docs: https://www.flowdock.com/api/flows
type FlowsService struct {
	client *Client
}

// Flow represents a Flowdock flow (room).
type Flow struct {
	Id                *string       `json:"id,omitempty"`
	Name              *string       `json:"name,omitempty"`
	ParameterizedName *string       `json:"parameterized_name,omitempty"`
	UnreadMentions    *int64        `json:"unread_mentions,omitempty"`
	Open              *bool         `json:"open,omitempty"`
	Disabled          *bool         `json:"disabled,omitempty"`
	Joined            *bool         `json:"joined,omitempty"`
	Url               *string       `json:"url,omitempty"`
	WebUrl            *string       `json:"web_url,omitempty"`
	JoinUrl           *string       `json:"join_url,omitempty"`
	AccessMode        *string       `json:"access_mode,omitempty"`
	Organization      *Organization `json:"organization,omitempty"`
	Users             *[]User       `json:"users,omitempty"`
}

// FlowsListOptions specifies the optional parameters to the FlowsService.List
// method.
type FlowsListOptions struct {
	// User a boolean value (1/0) that controls whether a list of users should
	// be included with each flow.
	User bool `url:"user,omitempty"`
}

// FlowsGetOptions specifies the optional parameters to the FlowsService.Get
// method.
type FlowsGetOptions struct {
	Id string `url:"id,omitempty"`
}

// FlowsCreateOptions specifies the optional parameters to the
// FlowsService.Create method.
type FlowsCreateOptions struct {
	Name string `url:"name"`
}

// Lists the flows that the authenticated user is a member of.
//
// Flowdock API docs: https://www.flowdock.com/api/flows
func (s *FlowsService) List(all bool, opt *FlowsListOptions) ([]Flow, *http.Response, error) {
	u := "flows"

	if all {
		u += "/all"
	}

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	flows := new([]Flow)
	resp, err := s.client.Do(req, flows)
	if err != nil {
		return nil, resp, err
	}

	return *flows, resp, err
}

// Get a single flow. Single flow information always includes the flow's user
// list. Otherwise, the data format is identical to what is returned with the
// list of flows.
//
// Flowdock API docs: https://www.flowdock.com/api/flows
func (s *FlowsService) Get(org, flowName string) (*Flow, *http.Response, error) {
	u := fmt.Sprintf("flows/%v/%v", org, flowName)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	flow := new(Flow)
	resp, err := s.client.Do(req, flow)
	if err != nil {
		return nil, resp, err
	}

	return flow, resp, err
}

// Get a single flow. Single flow information always includes the flow's user
// list. Otherwise, the data format is identical to what is returned with the
// list of flows.
//
// Flowdock API docs: https://www.flowdock.com/api/flows
func (s *FlowsService) GetById(id string) (*Flow, *http.Response, error) {
	u := "flows/find"
	u, err := addOptions(u, FlowsGetOptions{Id: id})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	flow := new(Flow)
	resp, err := s.client.Do(req, flow)
	if err != nil {
		return nil, resp, err
	}

	return flow, resp, err
}

// Create a flow for the specified organization
//
// Flowdock API docs: https://www.flowdock.com/api/flows
func (s *FlowsService) Create(orgName string, opt *FlowsCreateOptions) (*Flow, *http.Response, error) {
	u := fmt.Sprintf("flows/%v", orgName)

	u, err := addOptions(u, opt)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	flow := new(Flow)
	resp, err := s.client.Do(req, flow)
	if err != nil {
		return nil, resp, err
	}

	return flow, resp, err
}

// Update a flow.
//
// Flowdock API docs: https://www.flowdock.com/api/flows
func (s *FlowsService) Update(orgName, flowName string, flow *Flow) (*Flow, *http.Response, error) {
	u := fmt.Sprintf("flows/%v/%v", orgName, flowName)
	req, err := s.client.NewRequest("PUT", u, flow)
	if err != nil {
		return nil, nil, err
	}

	flow = new(Flow)
	resp, err := s.client.Do(req, flow)
	if err != nil {
		return nil, resp, err
	}

	return flow, resp, err
}
