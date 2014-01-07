package flowdock

import (
	"net/http"
	"encoding/json"
)

// MessagesService handles communication with the messages related methods of
// the Flowdock API.
//
// Flowdock API docs: https://www.flowdock.com/api/messages
type MessagesService struct {
	client *Client
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
	UserID           *string          `json:"user,omitempty"`
	Event            *string          `json:"event,omitempty"`
	RawContent       *json.RawMessage `json:"content,omitempty"`
	MessageID        *string          `json:"message,omitempty"`
	Tags             *[]string        `json:"tags,omitempty"`
	UUID             *string          `json:"uuid,omitempty"`
	ExternalUserName *string          `json:"external_user_name,omitempty"`
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
	}
	return content
}

// MessagesCreateOptions specifies the optional parameters to the
// MessageService.Create method.
type MessagesCreateOptions struct {
	FlowID           string   `url:"flow,omitempty"`
	MessageID        int      `url:"message,omitempty"`
	Event            string   `url:"event,omitempty"`
	Content          string   `url:"content,omitempty"`
	Tags             []string `url:"tags,omitempty"`
	UUID             string   `url:"uuid,omitempty"`
	ExternalUserName string   `url:"external_user_name,omitempty"`
}

// Content should be implemented by any value that is parsed into
// Message.RawContent. Its API will likly expand as more Message types are
// implemented.
type Content interface {
	String() string
}

// MessageContent represents a Message's Content when Message.Event is "message"
type MessageContent string

func (c *MessageContent) String() string {
	return string(*c)
}

// CommentContent represents a Message's Content when Message.Event is "comment"
type CommentContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
}

func (c *CommentContent) String() string {
	return *c.Text
}
