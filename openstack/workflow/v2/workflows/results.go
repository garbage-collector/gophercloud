package workflows

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Post operations. Call its Extract method to interpret it as a Workflow.
type CreateResult struct {
	commonResult
}

// GetResult is the response of Get operations. Call its Extract method to interpret it as a Workflow.
type GetResult struct {
	commonResult
}

// DeleteResult is the result from a Delete operation. Call its ExtractErr method to determine the success of the call.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract helps to get a Workflow struct from a Get function.
func (r commonResult) Extract() (*Workflow, error) {
	c := Workflow{}
	err := r.ExtractInto(&c)
	return &c, err
}

// Extract helps to get created Workflow struct from a Create function.
func (r CreateResult) Extract() ([]Workflow, error) {
	var s struct {
		Workflows []Workflow `json:"workflows"`
	}
	err := r.ExtractInto(&s)
	return s.Workflows, err
}

// Workflow represents a workflow execution on OpenStack mistral API.
type Workflow struct {
	// ID is the workflow's unique ID.
	ID string `json:"id"`

	// Definition is the workflow definition in Mistral v2 DSL.
	Definition string `json:"definition"`

	// Name is the name of the workflow.
	Name string `json:"name"`

	// Namespace is the namespace of the workflow.
	Namespace string `json:"namespace"`

	// Input represents the needed input to execute the workflow.
	// This parameter is a list of each input, comma separated.
	Input string `json:"input"`

	// ProjectID is the project id owner of the workflow.
	ProjectID string `json:"project_id"`

	// Scope is the scope of the workflow.
	// Values can be "private" or "public".
	Scope string `json:"scope"`
}

// WorkflowPage contains a single page of all workflows from a List call.
type WorkflowPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks if an WorkflowPage contains any results.
func (e WorkflowPage) IsEmpty() (bool, error) {
	exec, err := ExtractWorkflows(e)
	return len(exec) == 0, err
}

// NextPageURL finds the next page URL in a page in order to navigate to the next page of results.
func (e WorkflowPage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}

	err := e.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	return s.Next, nil
}

// ExtractWorkflows get the list of cron triggers from a page acquired from the List call.
func ExtractWorkflows(r pagination.Page) ([]Workflow, error) {
	var s struct {
		Workflows []Workflow `json:"workflows"`
	}
	err := (r.(WorkflowPage)).ExtractInto(&s)
	return s.Workflows, err
}
