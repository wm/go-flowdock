package flowdock

import (
	"testing"
	"net/http"
	"fmt"
	"reflect"
)

// lists users
// list users in a flow
// get a user by id
// updated a user by id
// add a user by id to a flow

func TestUsersService_List(t *testing.T) {
	nickOne := "wm"
	nickTwo := "om"

	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"1", "nick": "wm"}, {"id":"2", "nick": "om"}]`)
	})

	users, _, err := client.Users.List()
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	want := []User{{ID: &idOne, Nick: &nickOne}, {ID: &idTwo, Nick: &nickTwo}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.List returned %+v, want %+v", users, want)
	}
}
