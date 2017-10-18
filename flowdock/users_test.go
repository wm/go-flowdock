package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var (
	userId1 int = 1
	userId2 int = 2
)

func TestUsersService_All(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	users, _, err := client.Users.All()
	if err != nil {
		t.Errorf("Users.All returned error: %v", err)
	}

	want := []User{{Id: &userId1}, {Id: &userId2}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.All returned %+v, want %+v", users, want)
	}
}

func TestUsersService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/flows/orgname/flowname/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	users, _, err := client.Users.List("orgname", "flowname")
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	want := []User{{Id: &userId1}, {Id: &userId2}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Users.List returned %+v, want %+v", users, want)
	}
}

func TestUsersService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	user, _, err := client.Users.Get(userId1)
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := User{Id: &userId1}
	if !reflect.DeepEqual(user.Id, want.Id) {
		t.Errorf("Users.Get returned %+v, want %+v", user.Id, want.Id)
	}
}

func TestUsersService_Update(t *testing.T) {
	setup()
	defer teardown()

	nick := "new-nick"

	mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"id":1, "nick":"new-nick"}`)
	})

	opts := &UserUpdateOptions{
		Nick: "new-nick",
	}
	user, _, err := client.Users.Update(userId1, opts)
	if err != nil {
		t.Errorf("Users.Update returned error: %v", err)
	}

	want := User{Nick: &nick}
	if !reflect.DeepEqual(user.Nick, want.Nick) {
		t.Errorf("Users.Update returned %+v, want %+v", user.Nick, want.Nick)
	}
}
