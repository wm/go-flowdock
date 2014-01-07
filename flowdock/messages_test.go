package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMessagesService_Create_message(t *testing.T) {
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

	contentStr := "Howdy-Doo @Jackie #awesome"
	if !reflect.DeepEqual(message.ID, message.ID) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.ID, message.ID)
	}
	if !reflect.DeepEqual(message.Event, message.Event) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Event, message.Event)
	}
	if !reflect.DeepEqual(message.Content(), contentStr) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Content(), contentStr)
	}
}

func TestMessagesService_Create_comment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"event": "message",
			"content": "Howdy-Doo @Jackie #awesome",
		})
		fmt.Fprint(w, `{
			"event": "message",
			"content":{ "title":"Title of parent", "text":"This is a comment" }
		}`)
	})

	opt := MessagesCreateOptions{
		Event: "message",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := client.Messages.Create(&opt)
	if err != nil {
		t.Errorf("Messages.Create returned error: %v", err)
	}

	contentStr := "Howdy-Doo @Jackie #awesome"
	if !reflect.DeepEqual(message.ID, message.ID) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.ID, message.ID)
	}
	if !reflect.DeepEqual(message.Event, message.Event) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Event, message.Event)
	}
	if !reflect.DeepEqual(message.Content(), contentStr) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Content(), contentStr)
	}
}
