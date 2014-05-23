package flowdock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var (
	organizationId1 int = 1
	organizationId2 int = 2
)

func TestOrganizationsService_All(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	organizations, _, err := client.Organizations.All()
	if err != nil {
		t.Errorf("Organizations.All returned error: %v", err)
	}

	want := []Organization{{Id: &organizationId1}, {Id: &organizationId2}}
	if !reflect.DeepEqual(organizations, want) {
		t.Errorf("Organizations.All returned %+v, want %+v", organizations, want)
	}
}

func TestOrganizationsService_GetByParameterizedName(t *testing.T) {
	setup()
	defer teardown()

	name := "parameterizedorgname"

	mux.HandleFunc("/organizations/parameterizedorgname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"parameterized_name":"parameterizedorgname"}`)
	})

	organization, _, err := client.Organizations.GetByParameterizedName(name)
	if err != nil {
		t.Errorf("Organizations.GetByParameterizedName returned error: %v", err)
	}

	want := Organization{ParameterizedName: &name}
	if !reflect.DeepEqual(organization.ParameterizedName, want.ParameterizedName) {
		t.Errorf("Organizations.GetByParameterizedName returned %+v, want %+v", organization.ParameterizedName, want.ParameterizedName)
	}
}

func TestOrganizationsService_GetById(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations/find?id=1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	organization, _, err := client.Organizations.GetById(organizationId1)
	if err != nil {
		t.Errorf("Organizations.GetById returned error: %v", err)
	}

	want := Organization{Id: &organizationId1}
	if !reflect.DeepEqual(organization.Id, want.Id) {
		t.Errorf("Organizations.GetById returned %+v, want %+v", organization.Id, want.Id)
	}
}

func TestOrganizationsService_Update(t *testing.T) {
	setup()
	defer teardown()

	name := "new-name"

	mux.HandleFunc("/organizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"id":1, "name":"new-name"}`)
	})

	opts := &OrganizationUpdateOptions{
		Name: name,
	}
	organization, _, err := client.Organizations.Update(organizationId1, opts)
	if err != nil {
		t.Errorf("Organizations.Update returned error: %v", err)
	}

	want := Organization{Name: &name}
	if !reflect.DeepEqual(organization.Name, want.Name) {
		t.Errorf("Organizations.Update returned %+v, want %+v", organization.Name, want.Name)
	}
}
