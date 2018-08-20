package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/workflow/v2/executions"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListExecutions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/executions", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `{
				"executions": [
					{
						"created_at": "1970-01-01T00:00:00.000000",
						"description": "this is a description",
						"id": "1",
						"input": "{}",
						"output": "{}",
						"params": "{}",
						"project_id": "p1",
						"state": "SUCCESS",
						"updated_at": "1970-01-01T00:00:00.000000",
						"workflow_id": "w1",
						"workflow_name": "flow",
						"workflow_namespace": "some_namespace"
					}
				],
				"next": "%s/executions?marker=1"
			}`, th.Server.URL)
		case "1":
			fmt.Fprintf(w, `{ "executions": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	pages := 0
	// Get all executions
	err := executions.List(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := executions.ExtractExecutions(page)
		if err != nil {
			return false, err
		}

		expected := []executions.Execution{
			{ID: "1", Description: "this is a description", Input: "{}", Output: "{}", Params: "{}", ProjectID: "p1", State: "SUCCESS", WorkflowID: "w1", WorkflowName: "flow", WorkflowNamespace: "some_namespace"},
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

func TestGetExecution(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/executions/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"execution": {
					"created_at": "1970-01-01T00:00:00.000000",
					"description": "this is a description",
					"id": "1",
					"input": "{}",
					"output": "{}",
					"params": "{}",
					"project_id": "p1",
					"state": "SUCCESS",
					"updated_at": "1970-01-01T00:00:00.000000",
					"workflow_id": "w1",
					"workflow_name": "flow",
					"workflow_namespace": "some_namespace"
				}
			}
		`)
	})

	actual, err := executions.Get(fake.ServiceClient(), "1").Extract()
	if err != nil {
		t.Fatalf("Unable to get execution: %v", err)
	}

	expected := &executions.Execution{
		ID:                "1",
		Description:       "this is a description",
		Input:             "{}",
		Output:            "{}",
		Params:            "{}",
		ProjectID:         "p1",
		State:             "SUCCESS",
		WorkflowID:        "w1",
		WorkflowName:      "flow",
		WorkflowNamespace: "some_namespace",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestCreateExecution(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/executions", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"execution": {
					"created_at": "1970-01-01T00:00:00.000000",
					"description": "this is a description",
					"id": "1",
					"input": "{}",
					"output": "{}",
					"params": "{}",
					"project_id": "p1",
					"state": "SUCCESS",
					"updated_at": "1970-01-01T00:00:00.000000",
					"workflow_id": "w1",
					"workflow_name": "flow",
					"workflow_namespace": "some_namespace"
				}
			}
		`)
	})

	opts := &executions.CreateOpts{
		WorkflowID:  "w1",
		Input:       "{}",
		Params:      "{}",
		Description: "this is a description",
	}

	actual, err := executions.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create execution: %v", err)
	}

	expected := &executions.Execution{
		ID:                "1",
		Description:       "this is a description",
		Input:             "{}",
		Output:            "{}",
		Params:            "{}",
		ProjectID:         "p1",
		State:             "SUCCESS",
		WorkflowID:        "w1",
		WorkflowName:      "flow",
		WorkflowNamespace: "some_namespace",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteExecution(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/executions/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := executions.Delete(fake.ServiceClient(), "1")
	th.AssertNoErr(t, res.Err)
}
