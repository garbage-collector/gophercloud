package executions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Post operations. Call its Extract method to interpret it as an Execution.
type CreateResult struct {
	commonResult
}

// GetResult is the response of Get operations. Call its Extract method to interpret it as an Execution.
type GetResult struct {
	commonResult
}

// DeleteResult is the result from a Delete operation. Call its ExtractErr method to determine the success of the call.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract helps to get an Execution struct from a Get or a Create function.
func (r commonResult) Extract() (*Execution, error) {
	e := Execution{}
	err := r.ExtractInto(&e)
	return &e, err
}

// Execution represents a workflow execution on OpenStack mistral API.
type Execution struct {
	// ID is the execution's unique ID.
	ID string `json:"id"`

	// Description is the description of the execution.
	Description string `json:"description"`

	// Input contains the workflow input values in a JSON stringified object.
	Input string `json:"input"`

	// Ouput contains the workflow output values in a JSON stringified object.
	Output string `json:"output"`

	// Params contains workflow type specific parameters in a JSON stringified object.
	Params string `json:"params"`

	// ProjectID is the project id owner of the execution.
	ProjectID string `json:"project_id"`

	// State is the current state of the execution. State can be one of: IDLE, RUNNING, SUCCESS, ERROR, PAUSED.
	State string `json:"state"`

	// WorkflowID is the ID of the workflow linked to the execution.
	WorkflowID string `json:"workflow_id"`

	// WorkflowName is the name of the workflow linked to the execution.
	WorkflowName string `json:"workflow_name"`

	// WorkflowNamespace is the namespace of the workflow linked to the execution.
	WorkflowNamespace string `json:"workflow_namespace"`
}

// ExecutionPage contains a single page of all executions from a List call.
type ExecutionPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks if an ExecutionPage contains any results.
func (e ExecutionPage) IsEmpty() (bool, error) {
	exec, err := ExtractExecutions(e)
	return len(exec) == 0, err
}

// NextPageURL finds the next page URL in a page in order to navigate to the next page of results.
func (e ExecutionPage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}

	err := e.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	return s.Next, nil
}

// ExtractExecutions get the list of executions from a page acquired from the List call.
func ExtractExecutions(r pagination.Page) ([]Execution, error) {
	var s struct {
		Executions []Execution `json:"executions"`
	}
	err := (r.(ExecutionPage)).ExtractInto(&s)
	return s.Executions, err
}
