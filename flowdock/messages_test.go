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
	})

	opt := MessagesCreateOptions{
		Event: "message",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := client.Messages.Create(&opt)
	if err != nil {
		t.Errorf("Messages.Create returned error: %v", err)
	}

	if !reflect.DeepEqual(*message.Event, opt.Event) {
		t.Errorf("Messages.Create returned %+v, want %+v", *message.Event, opt.Event)
	}

	if !reflect.DeepEqual(message.Content().String(), opt.Content) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Content(), opt.Content)
	}
}

func TestMessagesService_Create_comment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"event": "comment",
			"content": "This is a comment",
		})
		fmt.Fprint(w, `{
			"event": "comment",
			"content":{ "title":"Title of parent", "text":"This is a comment" }
		}`)
	})

	opt := MessagesCreateOptions{
		Event: "comment",
		Content: "This is a comment",
	}
	message, _, err := client.Messages.CreateComment(&opt)
	if err != nil {
		t.Errorf("Messages.CreateComment returned error: %v", err)
	}

	if !reflect.DeepEqual(message.Event, message.Event) {
		t.Errorf("Messages.Create returned %+v, want %+v", message.Event, message.Event)
	}

	content        := CommentContent{Title: String("Title of parent"), Text: String("This is a comment")}
	messageContent := message.Content()
	if !reflect.DeepEqual(messageContent, &content) {
		t.Errorf("Messages.Create returned %+v, want %+v", messageContent, &content)
	}
}
