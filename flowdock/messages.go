package flowdock

import (
	"net/http"
	"encoding/json"
	"fmt"
)

// MessagesService handles communication with the messages related methods of
// the Flowdock API.
//
// Flowdock API docs: https://www.flowdock.com/api/messages
type MessagesService struct {
	client *Client
}

// MessagesListOptions specifies the optional parameters to the
// MessageService.List method.
type MessagesListOptions struct {
	Event            string   `url:"event,omitempty"`
	Limit            int      `url:"limit,omitempty"`
	SinceId          int      `url:"since_id,omitempty"`
	UntilId          int      `url:"until_id,omitempty"`
	Tags             []string `url:"tags,comma,omitempty"`
	TagMode          string   `url:"tag_mode,omitempty"`
	Search           string   `url:"search,omitempty"`
}

// Lists the messages for the given flow.
//
// Flowdock API docs: https://www.flowdock.com/api/messages
func (s *MessagesService) List(org, flow string, opt *MessagesListOptions) ([]Message, *http.Response, error) {
	u := fmt.Sprintf("flows/%v/%v/messages", org, flow)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	messages := new([]Message)
	resp, err := s.client.Do(req, messages)
	if err != nil {
		return nil, resp, err
	}

	return *messages, resp, err
}

// MessagesCreateOptions specifies the optional parameters to the
// MessageService.Create method.
type MessagesCreateOptions struct {
	FlowID           string   `url:"flow,omitempty"`
	MessageID        int      `url:"message,omitempty"`
	Event            string   `url:"event,omitempty"`
	Content          string   `url:"content,omitempty"`
	Tags             []string `url:"tags,comma,omitempty"`
	UUID             string   `url:"uuid,omitempty"`
	ExternalUserName string   `url:"external_user_name,omitempty"`
	Subject          string   `url:"subject,omitempty"`
	FromAddress      string   `url:"from_address,omitempty"`
	Source           string   `url:"source,omitempty"`
}

// Create a comment for the specified organization
//
// Flowdock API docs: https://www.flowdock.com/api/messages
func (s *MessagesService) CreateComment(opt *MessagesCreateOptions) (*Message, *http.Response, error) {
	u := "comments"

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

// Message represents a Flowdock chat message.
type Message struct {
	ID               *int             `json:"id,omitempty"`
	FlowID           *string          `json:"flow,omitempty"`
	Sent             *Time            `json:"sent,omitempty"`
	UserID           *int             `json:"user,omitempty"`
	Event            *string          `json:"event,omitempty"`
	RawContent       *json.RawMessage `json:"content,omitempty"`
	MessageID        *int             `json:"message,omitempty"`
	Tags             *[]string        `json:"tags,omitempty"`
	UUID             *string          `json:"uuid,omitempty"`
	ExternalUserName *string          `json:"external_user_name,omitempty"`
	App              *string          `json:"app,omitempty"` // deprecated
}

// Return the content of a Message
//
// It can be a MessageContent, CommentContent, etc. Depends on the Event
func (m *Message) Content() (content Content) {
	switch *m.Event {
	case "message":
		content = new(MessageContent)
		if err := json.Unmarshal([]byte(*m.RawContent), &content); err != nil {
			panic(err.Error())
		}
	case "comment":
		content = &CommentContent{}
		if err := json.Unmarshal([]byte(*m.RawContent), &content); err != nil {
			panic(err.Error())
		}
	// default:
	// 	messageContent := MessageContent(string(*m.RawContent))
	// 	content = &messageContent
	}

	return content
}

// Content should be implemented by any value that is parsed into
// Message.RawContent. Its API will likly expand as more Message types are
// implemented.
type Content interface {
	String() string
}

// MessageContent represents a Message's Content when Message.Event is "message"
type MessageContent string

// Return the string version of a MessageContent
//
func (c *MessageContent) String() string {
	return string(*c)
}

// CommentContent represents a Message's Content when Message.Event is "comment"
type CommentContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
}

// Return the string version of a CommentContent
//
// It returns the *CommentContent.Text
func (c *CommentContent) String() string {
	return *c.Text
}
