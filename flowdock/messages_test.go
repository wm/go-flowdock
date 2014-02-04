package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMessagesService_Stream(t *testing.T) {
	setup()
	defer teardown()
	more := make(chan bool, 1)

	mux.HandleFunc("/flows/org/flow", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"access_token": "token"})
		w.Header().Set("Content-Type", "text/event-stream")

		var id int

		// send a message for each 'more'
		for {
			if !<-more {
				break
			}

			fmt.Fprintf(w, "id: %d\ndata: {\"event\":\"message\",\"content\":\"message %d\"}\n\n", id, id)
			w.(responseWriter).Flush()
			id++
		}
	})
	defer close(more)

	stream, _, err := client.Messages.Stream("token", "org", "flow")
	more <- true // tell test server to send a message

	if err != nil {
		t.Errorf("Messages.Stream returned error: %v", err)
	}

	msg := <-stream

	if msg.Content().String() != "message 0" {
		t.Fatalf("expected message 0, got %v", msg.Content())
	}
}

func TestMessagesService_List(t *testing.T) {
	setup()
	defer teardown()
	var idOne      = 3816534
	var eventOne   = "message"
	var content    = []string{"Hello NYC", "Hello World"}
	var idTwo      = 45590
	var eventTwo   = "message"

	mux.HandleFunc("/flows/org/flow/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
		  {
			"app":"chat",
			"sent":1317397485508,
			"uuid":"odHapx1VWp7WTrdQ",
			"tags":[],
			"flow": "deadbeefdeadbeef",
			"id":3816534,
			"event":"message",
			"content": "Hello NYC",
			"attachments": [],
			"user":"18"
		  },
		  {
			"app": "chat",
			"event": "message",
			"tags": [],
			"uuid": "4W_LQEybVaX-gJmi",
			"id": 45590,
			"flow": "deadbeefdeadbeef",
			"content": "Hello World",
			"sent": 1317715340213,
			"attachments": [],
			"user": "2"
		  }
		]`)
		fmt.Fprint(w, `[{"id":"1"}, {"id":"2"}]`)
	})

	messages, _, err := client.Messages.List("org", "flow", nil)
	if err != nil {
		t.Errorf("Messages.List returned error: %v", err)
	}

	want := []Message{
		{
			ID: &idOne,
			Event: &eventOne,
		},
		{
			ID: &idTwo,
			Event: &eventTwo,
		},
	}

	for i, msg := range messages {
		if *msg.ID != *want[i].ID {
			t.Errorf("Messages.List returned %+v, want %+v", *msg.ID, *want[i].ID)
		}
		if *msg.Event != *want[i].Event {
			t.Errorf("Messages.List returned %+v, want %+v", *msg.Event, *want[i].Event)
		}
		if msg.Content().String() != content[i] {
			t.Errorf("Messages.List returned %+v, want %+v", msg.Content(), content[i])
		}
	}
}

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

	title          := "Title of parent"
	text           := "This is a comment"
	content        := CommentContent{Title: &title, Text: &text}
	messageContent := message.Content()
	if !reflect.DeepEqual(messageContent, &content) {
		t.Errorf("Messages.Create returned %+v, want %+v", messageContent, &content)
	}
}

func TestCommentContent_String(t *testing.T) {
	title   := "Title of parent"
	text    := "This is a comment"
	content := CommentContent{Title: &title, Text: &text}

	want    := "This is a comment"
	if (*content.Text != want) {
		t.Errorf("Messages.Create returned %+v, want %+v", *content.Text, want)
	}
}
