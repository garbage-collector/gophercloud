package v2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/workflow/v2/crontriggers"
	"github.com/gophercloud/gophercloud/openstack/workflow/v2/executions"
	"github.com/gophercloud/gophercloud/openstack/workflow/v2/workflows"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateWorkflow creates a workflow on Mistral API.
// The created workflow is a dummy workflow that performs a simple echo.
func CreateWorkflow(t *testing.T, client *gophercloud.ServiceClient) (*workflows.Workflow, error) {
	workflowName := tools.RandomString("workflow_create_vm_", 5)

	definition := `---
version: '2.0'

` + workflowName + `:
  description: Simple workflow example
  type: direct

  tasks:
    test:
      action: std.echo output="Hello World!"`

	t.Logf("Attempting to create workflow: %s", workflowName)

	opts := &workflows.CreateOpts{
		Namespace:  "some-namespace",
		Scope:      "private",
		Definition: strings.NewReader(definition),
	}
	workflowList, err := workflows.Create(client, opts).Extract()
	if err != nil {
		return nil, err
	}
	th.AssertEquals(t, 1, len(workflowList))

	workflow := workflowList[0]

	t.Logf("Workflow created: %s", workflowName)

	th.AssertEquals(t, workflowName, workflow.Name)

	return &workflow, nil
}

// DeleteWorkflow deletes the given workflow.
func DeleteWorkflow(t *testing.T, client *gophercloud.ServiceClient, workflow *workflows.Workflow) {
	err := workflows.Delete(client, workflow.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete workflows %s: %v", workflow.Name, err)
	}

	t.Logf("Deleted workflow: %s", workflow.Name)
}

// CreateExecution creates an execution for the given workflow.
// This method waits the success of the execution.
func CreateExecution(t *testing.T, client *gophercloud.ServiceClient, workflow *workflows.Workflow) (*executions.Execution, error) {
	executionDescription := tools.RandomString("execution_", 5)

	t.Logf("Attempting to create execution: %s", executionDescription)
	createOpts := executions.CreateOpts{
		WorkflowID:        workflow.ID,
		WorkflowNamespace: workflow.Namespace,
		Description:       executionDescription,
		Input:             "{}",
		Params:            "{}",
	}
	execution, err := executions.Create(client, createOpts).Extract()
	if err != nil {
		return execution, err
	}

	t.Logf("Execution created: %s", executionDescription)

	th.AssertEquals(t, execution.Description, executionDescription)

	t.Logf("Wait for execution status SUCCESS: %s", executionDescription)
	th.AssertNoErr(t, tools.WaitFor(func() (bool, error) {
		latest, err := executions.Get(client, execution.ID).Extract()
		if err != nil {
			return false, err
		}

		if latest.State == "SUCCESS" {
			execution = latest
			return true, nil
		}

		if latest.State == "ERROR" {
			return false, fmt.Errorf("Execution in ERROR state")
		}

		return false, nil
	}))
	t.Logf("Execution success: %s", executionDescription)

	return execution, nil
}

// DeleteExecution deletes an execution.
func DeleteExecution(t *testing.T, client *gophercloud.ServiceClient, execution *executions.Execution) {
	err := executions.Delete(client, execution.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete executions %s: %v", execution.Description, err)
	}

	t.Logf("Deleted executions: %s", execution.Description)
}

// CreateCronTrigger creates a cron trigger for the given workflow.
func CreateCronTrigger(t *testing.T, client *gophercloud.ServiceClient, workflow *workflows.Workflow) (*crontriggers.CronTrigger, error) {
	crontriggerName := tools.RandomString("crontrigger_", 5)

	t.Logf("Attempting to create cron trigger: %s", crontriggerName)
	createOpts := crontriggers.CreateOpts{
		WorkflowID:     workflow.ID,
		Name:           crontriggerName,
		Pattern:        "0 0 1 1 *",
		WorkflowInput:  "{}",
		WorkflowParams: "{}",
	}
	crontrigger, err := crontriggers.Create(client, createOpts).Extract()
	if err != nil {
		return crontrigger, err
	}

	t.Logf("Cron trigger created: %s", crontriggerName)

	th.AssertEquals(t, crontrigger.Name, crontriggerName)

	return crontrigger, nil
}

// DeleteCronTrigger deletes a cron trigger.
func DeleteCronTrigger(t *testing.T, client *gophercloud.ServiceClient, crontrigger *crontriggers.CronTrigger) {
	err := crontriggers.Delete(client, crontrigger.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete cron trigger %s: %v", crontrigger.Name, err)
	}

	t.Logf("Deleted crontrigger: %s", crontrigger.Name)
}
