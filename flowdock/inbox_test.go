package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestInboxService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/messages/team_inbox/xxx", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"subject": "a subject",
			"content": "Howdy-Doo @Jackie #awesome",
		})
		fmt.Fprint(w, `{}`)
	})

	opt := InboxCreateOptions{
		Subject: "a subject",
		Content: "Howdy-Doo @Jackie #awesome",
	}
	message, _, err := client.Inbox.Create("xxx", &opt)
	if err != nil {
		t.Errorf("Messages.Create returned error: %v", err)
	}

	want := new(Message)
	if !reflect.DeepEqual(message, want) {
		t.Errorf("Messages.Create returned \n%+v \nwant \n%+v", message, want)
	}
}
