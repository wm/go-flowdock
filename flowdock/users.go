package flowdock

import (
	"fmt"
	"net/http"
)

type UserUpdateOptions struct {
	Nick  string `json:"nick,omitempty"`
	Email string `json:"email,omitempty"`
}

type UsersService struct {
	client *Client
}

// All users visible to the authenticated user.
//
// Flowdock API docs: https://www.flowdock.com/api/users
func (s *UsersService) All() ([]User, *http.Response, error) {
	u := "users"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	users := new([]User)
	resp, err := s.client.Do(req, users)
	if err != nil {
		return nil, resp, err
	}

	return *users, resp, err
}

// List the users inside a flow.
//
// Flowdock API docs: https://www.flowdock.com/api/users
func (s *UsersService) List(org, flow string) ([]User, *http.Response, error) {
	u := fmt.Sprintf("flows/%v/%v/users", org, flow)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	users := new([]User)
	resp, err := s.client.Do(req, users)
	if err != nil {
		return nil, resp, err
	}

	return *users, resp, err
}

// Get a user by their id.
//
// Flowdock API docs: https://www.flowdock.com/api/users
func (s *UsersService) Get(id int) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%v", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

// Update a user by their id.
//
// Flowdock API docs: https://www.flowdock.com/api/users
func (s *UsersService) Update(id int, opt *UserUpdateOptions) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%v", id)

	u, err := addOptions(u, opt)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

type User struct {
	Id           *int    `json:"id,omitempty"`
	Nick         *string `json:"nick,omitempty"`
	Name         *string `json:"name,omitempty"`
	Email        *string `json:"email,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
	Status       *string `json:"status,omitempty"`
	Disabled     *bool   `json:"disabled,omitempty"`
	LastActivity *Time   `json:"last_activity,omitempty"`
	LastPing     *Time   `json:"last_ping,omitempty"`
}
