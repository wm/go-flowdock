package flowdock

import (
	"fmt"
	"net/http"
	"net/url"
)

type OrganizationUpdateOptions struct {
	Name string `json:"name,omitempty"`
}

type OrganizationsService struct {
	client *Client
}

// All organizations authenticated user belongs to.
//
// Flowdock API docs: https://www.flowdock.com/api/organizations
func (s *OrganizationsService) All() ([]Organization, *http.Response, error) {
	u := "organizations"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	organizations := new([]Organization)
	resp, err := s.client.Do(req, organizations)
	if err != nil {
		return nil, resp, err
	}

	return *organizations, resp, err
}

// GetByParameterizedName fetches an organization by it's parameterized_name.
//
// Flowdock API docs: https://www.flowdock.com/api/organizations
func (s *OrganizationsService) GetByParameterizedName(name string) (*Organization, *http.Response, error) {
	u := fmt.Sprintf("organizations/%v", name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	organization := new(Organization)
	resp, err := s.client.Do(req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// GetById fetches an organization by it's id.
//
// Flowdock API docs: https://www.flowdock.com/api/organizations
func (s *OrganizationsService) GetById(id int) (*Organization, *http.Response, error) {
	u := fmt.Sprintf("organizations/find?id=%v", id)

	req, err := s.client.NewRequest("GET", url.QueryEscape(u), nil)
	if err != nil {
		return nil, nil, err
	}

	organization := new(Organization)
	resp, err := s.client.Do(req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Update an organization by id.
//
// Flowdock API docs: https://www.flowdock.com/api/organizations
func (s *OrganizationsService) Update(id int, opt *OrganizationUpdateOptions) (*Organization, *http.Response, error) {
	u := fmt.Sprintf("organizations/%v", id)

	u, err := addOptions(u, opt)
	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}

	organization := new(Organization)
	resp, err := s.client.Do(req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Organization represents a Flowdock organization to which members belong
type Organization struct {
	Id                *int    `json:"id,omitempty"`
	Name              *string `json:"name,omitempty"`
	ParameterizedName *string `json:"parameterized_name,omitempty"`
	UserLimit         *int64  `json:"user_limit,omitempty"`
	UserCount         *int64  `json:"user_count,omitempty"`
	Active            *bool   `json:"active,omitempty"`
	Url               *string `json:"url,omitempty"`
	Users             *[]User `json:"users"`
}
