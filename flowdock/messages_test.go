package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"encoding/json"
)

func TestMessagesService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"event": "message",
			"content": "Howdy-Doo @Jackie #awesome",
		})
		fmt.Fprint(w, `{
			"event": "message",
			"content": "Howdy-Doo @Jackie #awesome"
		}`)
			// "content":{ "title":"Title of parent", "text":"This is a comment" }
	})

	opt := MessagesCreateOptions{
		Event: "message",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := client.Messages.Create(&opt)
	if err != nil {
		t.Errorf("Messages.Create returned error: %v", err)
	}

	rawMessage := json.RawMessage("Howdy-Doo @Jackie #awesome")
	want := &Message{
		ID: message.ID,
		Event: String("message"),
		RawContent: &rawMessage,
	}
	if !reflect.DeepEqual(message.ID, want.ID) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.ID, want.ID)
	}
	if !reflect.DeepEqual(message.Event, want.Event) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Event, want.Event)
	}
	if !reflect.DeepEqual(message.Content(), want.Content()) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Content(), want.Content())
	}
}
