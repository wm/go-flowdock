package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
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
	})

	opt := MessagesCreateOptions{
		Event: "message",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := client.Messages.Create(&opt)
	if err != nil {
		t.Errorf("Messages.Create returned error: %v", err)
	}

	want := &Message{
		ID: message.ID,
		Event: String("message"),
		Content: String("Howdy-Doo @Jackie #awesome"),
	}
	if !reflect.DeepEqual(message, want) {
		t.Errorf("Messages.Create returned %+v, want %+v", message, want)
	}
}
