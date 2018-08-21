package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"

	"github.com/gophercloud/gophercloud/openstack/workflow/v2/crontriggers"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListCronTriggers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `{
				"cron_triggers": [
					{
						"created_at": "1970-01-01 00:00:00",
						"id": "1",
						"name": "trigger",
						"pattern": "* * * * *",
						"project_id": "p1",
						"remaining_executions": 42,
						"scope": "private",
						"updated_at": "1970-01-01 00:00:00",
						"workflow_id": "w1",
						"workflow_input": "{}",
						"workflow_name": "my_wf",
						"workflow_params": "{}"
					}
				],
				"next": "%s/cron_triggers?marker=1"
			}`, th.Server.URL)
		case "1":
			fmt.Fprintf(w, `{ "cron_triggers": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})

	pages := 0
	// Get all cron triggers
	err := crontriggers.List(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := crontriggers.ExtractCronTriggers(page)
		if err != nil {
			return false, err
		}

		expected := []crontriggers.CronTrigger{
			{ID: "1", Name: "trigger", Pattern: "* * * * *", ProjectID: "p1", RemainingExecutions: 42, Scope: "private", WorkflowID: "w1", WorkflowName: "my_wf", WorkflowInput: "{}", WorkflowParams: "{}", CreatedAt: gophercloud.JSONRFC3339ZNoTNoZ(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC))},
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

func TestGetCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "1970-01-01 00:00:00",
				"id": "1",
				"name": "trigger",
				"pattern": "* * * * *",
				"project_id": "p1",
				"remaining_executions": 42,
				"scope": "private",
				"updated_at": "1970-01-01 00:00:00",
				"workflow_id": "w1",
				"workflow_input": "{}",
				"workflow_name": "my_wf",
				"workflow_params": "{}"
			}
		`)
	})

	actual, err := crontriggers.Get(fake.ServiceClient(), "1").Extract()
	if err != nil {
		t.Fatalf("Unable to get cron trigger: %v", err)
	}

	expected := &crontriggers.CronTrigger{
		ID:                  "1",
		Name:                "trigger",
		Pattern:             "* * * * *",
		ProjectID:           "p1",
		RemainingExecutions: 42,
		Scope:               "private",
		WorkflowID:          "w1",
		WorkflowName:        "my_wf",
		WorkflowInput:       "{}",
		WorkflowParams:      "{}",
		CreatedAt:           gophercloud.JSONRFC3339ZNoTNoZ(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestCreateCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "1970-01-01 00:00:00",
				"id": "1",
				"name": "trigger",
				"pattern": "* * * * *",
				"project_id": "p1",
				"remaining_executions": 42,
				"scope": "private",
				"updated_at": "1970-01-01 00:00:00",
				"workflow_id": "w1",
				"workflow_input": "{}",
				"workflow_name": "my_wf",
				"workflow_params": "{}"
			}
		`)
	})

	opts := &crontriggers.CreateOpts{
		WorkflowID:     "w1",
		WorkflowInput:  "{}",
		WorkflowParams: "{}",
		Name:           "trigger",
	}

	actual, err := crontriggers.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create cron trigger: %v", err)
	}

	expected := &crontriggers.CronTrigger{
		ID:                  "1",
		Name:                "trigger",
		Pattern:             "* * * * *",
		ProjectID:           "p1",
		RemainingExecutions: 42,
		Scope:               "private",
		WorkflowID:          "w1",
		WorkflowName:        "my_wf",
		WorkflowInput:       "{}",
		WorkflowParams:      "{}",
		CreatedAt:           gophercloud.JSONRFC3339ZNoTNoZ(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := crontriggers.Delete(fake.ServiceClient(), "1")
	th.AssertNoErr(t, res.Err)
}
