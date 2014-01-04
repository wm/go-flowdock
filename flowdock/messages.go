package flowdock

import (
	"net/http"
)

// MessagesService handles communication with the messages related methods of
// the Flowdock API.
//
// Flowdock API docs: https://www.flowdock.com/api/messages
type MessagesService struct {
	client *Client
}

// Message represents a Flowdock chat message.
type Message struct {
	ID               *int       `json:"id,omitempty"`
	FlowID           *string    `json:"flow,omitempty"`
	Sent             *Time      `json:"sent,omitempty"`
	UserID           *string    `json:"user,omitempty"`
	Event            *string    `json:"event,omitempty"`
	Content          *string    `json:"content,omitempty"`
	MessageID        *string    `json:"message,omitempty"`
	Tags             *[]string  `json:"tags,omitempty"`
	UUID             *string    `json:"uuid,omitempty"`
	ExternalUserName *string    `json:"external_user_name,omitempty"`
}

// MessagesCreateOptions specifies the optional parameters to the
// MessageService.Create method.
type MessagesCreateOptions struct {
	FlowID           string   `url:"flow,omitempty"`
	Event            string   `url:"event,omitempty"`
	Content          string   `url:"content,omitempty"`
	MessageID        string   `url:"message,omitempty"`
	Tags             []string `url:"tags,omitempty"`
	UUID             string   `url:"uuid,omitempty"`
	ExternalUserName string   `url:"external_user_name,omitempty"`
}

// Create a message for the specified organization
//
// Flowdock API docs: https://www.flowdock.com/api/messages
func (s *MessagesService) Create(opt *MessagesCreateOptions) (*Message, *http.Response, error) {
	u := "messages"

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
