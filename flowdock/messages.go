package flowdock

import (
	"time"
)

// Message represents a Flowdock chat message.
type Message struct {
	ID        string
	FlowID    string `json:"flow"`
	Sent      time.Time
	UserID    string `json:"user"`
	Event     string
	Content   string
	MessageID string `json:"message"`
	Tags      []string
	UUID      string
}
