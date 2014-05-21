package flowdock

import (
	"fmt"
	"net/http"
)

// InboxService handles communication with the Team Inbox related methods of
// the Flowdock API.
//
// Flowdock API docs: https://flowdock.com/api/team-inbox
type InboxService struct {
	client *Client
}

// InboxCreateOptions specifies the optional parameters to the
// InboxCreate method.
type InboxCreateOptions struct {
	Source      string   `url:"source,omitempty"`
	FromAddress string   `url:"from_address,omitempty"`
	Subject     string   `url:"subject,omitempty"`
	Content     string   `url:"content,omitempty"`
	FromName    string   `url:"from_name,omitempty"`
	ReplyTo     string   `url:"reply_to,omitempty"`
	Project     string   `url:"project,omitempty"`
	Tags        []string `url:"tags,comma,omitempty"`
	Link        string   `url:"link,omitempty"`
}

// Create an Inbox mail message for the specified flow api token
//
// Flowdock API docs: https://www.flowdock.com/api/team-inbox
func (s *InboxService) Create(flowApiToken string, opt *InboxCreateOptions) (*Message, *http.Response, error) {
	u := fmt.Sprintf("v1/messages/team_inbox/%v", flowApiToken)

	u, err := addOptions(u, opt)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	message := new(Message)
	resp, err := s.client.Do(req, message)
	if err != nil {
		return nil, resp, err
	}

	return message, resp, err
}
