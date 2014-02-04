package flowdock

import (
	"net/http"
	"fmt"
)

// UsersService handles communication with the user related methods of the
// Flowdock API.
//
// Flowdock API docs: https://www.flowdock.com/api/users
type UsersService struct {
	client *Client
}

type User struct {
	ID            *int     `json:"id,omitempty"`
	Nick          *string  `json:"nick,omitempty"`
	Name          *string  `json:"name,omitempty"`
	Email         *string  `json:"email,omitempty"`
	Avatar        *string  `json:"avatar,omitempty"`
	Status        *string  `json:"status,omitempty"`
	Disabled      *bool    `json:"disabled,omitempty"`
	LastActivity  *Time    `json:"last_activity,omitempty"`
	LastPing      *Time    `json:"last_ping,omitempty"`
}

// List all users visible to the authenticated user, i.e. a combined set of
// users from all the organizations of the authenticated user. If the
// authenticated user is an admin in an organization, all of that
// organization's users are returned. Otherwise, only users that are in the
// same flows as the authenticated user are returned.
//
// Flowdock API docs: https://www.flowdock.com/api/users
func (s *UsersService) List() (*[]User, *http.Response, error) {
	u := "/users"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	users := new([]User)
	resp, err := s.client.Do(req, users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, err
}

// Flowdock API docs: https://www.flowdock.com/api/users
func (s *UsersService) Get(id int) (*User, *http.Response, error) {
	u := fmt.Sprintf("/users/%v", id)

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
