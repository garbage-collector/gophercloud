package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/workflow/v2/workflows"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListWorkflows(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/workflows", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `{
				"next": "%s/workflows?marker=1",
				"workflows": [
					{
						"created_at": "1970-01-01T00:00:00.000000",
						"definition": "Workflow Definition in Mistral DSL v2",
						"id": "1",
						"input": "param1, param2",
						"name": "flow",
						"namespace": "some-namespace",
						"project_id": "p1",
						"scope": "private",
						"updated_at": "1970-01-01T00:00:00.000000"
					}
				]
			}`, th.Server.URL)
		case "1":
			fmt.Fprintf(w, `{ "workflows": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	pages := 0
	// Get all workflows
	err := workflows.List(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := workflows.ExtractWorkflows(page)
		if err != nil {
			return false, err
		}

		expected := []workflows.Workflow{
			{ID: "1", Definition: "Workflow Definition in Mistral DSL v2", Name: "flow", Namespace: "some-namespace", Input: "param1, param2", ProjectID: "p1", Scope: "private"},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %#v, but was %#v", expected, actual)
		}

		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if pages != 1 {
		t.Errorf("Expected one page, got %d", pages)
	}
}

func TestGetWorkflow(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/workflows/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "1970-01-01T00:00:00.000000",
				"definition": "Workflow Definition in Mistral DSL v2",
				"id": "1",
				"input": "param1, param2",
				"name": "flow",
				"namespace": "some-namespace",
				"project_id": "p1",
				"scope": "private",
				"updated_at": "1970-01-01T00:00:00.000000"
			}
		`)
	})

	actual, err := workflows.Get(fake.ServiceClient(), "1").Extract()
	if err != nil {
		t.Fatalf("Unable to get workflow: %v", err)
	}

	expected := &workflows.Workflow{
		ID:         "1",
		Definition: "Workflow Definition in Mistral DSL v2",
		Name:       "flow",
		Namespace:  "some-namespace",
		Input:      "param1, param2",
		ProjectID:  "p1",
		Scope:      "private",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestCreateWorkflow(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	definition := `---
version: '2.0'

create_vm:
	description: Simple workflow example
	type: direct

	input:
	- vm_name
	- image_ref
	- flavor_ref
	output:
	vm_id: <% $.vm_id %>

	tasks:
	create_server:
		action: nova.servers_create name=<% $.vm_name %> image=<% $.image_ref %> flavor=<% $.flavor_ref %>
		publish:
		vm_id: <% task(create_server).result.id %>
		on-success:
		- wait_for_instance

	wait_for_instance:
		action: nova.servers_find id=<% $.vm_id %> status='ACTIVE'
		retry:
		delay: 5
		count: 15`

	th.Mux.HandleFunc("/workflows", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestFormValues(t, r, map[string]string{
			"namespace": "some-namespace",
			"scope":     "private",
		})
		th.TestBody(t, r, definition)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
			{
				"created_at": "1970-01-01T00:00:00.000000",
				"definition": "Workflow Definition in Mistral DSL v2",
				"id": "1",
				"input": "param1, param2",
				"name": "flow",
				"namespace": "some-namespace",
				"project_id": "p1",
				"scope": "private",
				"updated_at": "1970-01-01T00:00:00.000000"
			}
		`)
	})

	opts := &workflows.CreateOpts{
		Namespace:  "some-namespace",
		Scope:      "private",
		Definition: strings.NewReader(definition),
	}

	actual, err := workflows.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create workflow: %v", err)
	}

	expected := &workflows.Workflow{
		ID:         "1",
		Definition: "Workflow Definition in Mistral DSL v2",
		Name:       "flow",
		Namespace:  "some-namespace",
		Input:      "param1, param2",
		ProjectID:  "p1",
		Scope:      "private",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteWorkflow(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/workflows/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := workflows.Delete(fake.ServiceClient(), "1")
	th.AssertNoErr(t, res.Err)
}
