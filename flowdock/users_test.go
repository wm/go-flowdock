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
	idOne := 1
	idTwo := 2
	nickOne := "wm"
	nickTwo := "om"

	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1, "nick": "wm"}, {"id":2, "nick": "om"}]`)
	})

	users, _, err := client.Users.List()
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	want := []User{{ID: &idOne, Nick: &nickOne}, {ID: &idTwo, Nick: &nickTwo}}
	if !reflect.DeepEqual(users, &want) {
		t.Errorf("Users.List returned %+v, want %+v", users, want)
	}
}

func TestUsersService_Get(t *testing.T) {
	id := 2
	nickOne := "wm"

	setup()
	defer teardown()

	mux.HandleFunc("/users/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":2, "nick": "wm"}`)
	})

	user, _, err := client.Users.Get(id)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := User{ID: &id, Nick: &nickOne}
	if !reflect.DeepEqual(user, &want) {
		t.Errorf("Users.Get returned %+v, want %+v", user, want)
	}
}
