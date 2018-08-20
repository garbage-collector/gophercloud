package executions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extension to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToExecutionListQuery() (string, error)
}

// ListOpts filters the result returned by the List() function.
type ListOpts struct {
	// WorkflowName allows to filter by workflow name.
	WorkflowName string `q:"workflow_name"`
	// WorkflowID allows to filter by workflow id.
	WorkflowID string `q:"workflow_id"`
	// Description allows to filter by execution description.
	Description string `q:"description"`

	// State allows to filter by execution state.
	State string `q:"state"`

	// SortDir allows to select sort direction.
	// It can be "asc" or "desc" (default).
	SortDir string `q:"sort_dir"`

	// SortKey allows to sort by one of the execution attributes.
	SortKey string `q:"sort_key"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of
	// executions.
	Limit int `q:"limit"`
}

// ToExecutionListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToExecutionListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List performs a call to list executions.
// You may provide options to filter the executions.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToExecutionListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ExecutionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extension to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToExecutionCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters used to create an execution.
type CreateOpts struct {
	// WorkflowID is the unique id of the workflow.
	WorkflowID string `json:"workflow_id" required:"true"`

	// WorkflowNamespace is the namespace of the workflow.
	WorkflowNamespace string `json:"workflow_namespace,omitempty"`

	// Input is a JSON structure containing workflow input values, serialized as string.
	Input string `json:"input,omitempty"`

	// Params define workflow type specific parameters.
	Params string `json:"params,omitempty"`

	// Description is the description of the workflow execution.
	Description string `json:"description,omitempty"`
}

// ToExecutionCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToExecutionCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "execution")
}

// Create requests the creation of a new execution.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToExecutionCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})

	return
}

// Get retrieves details of a single execution.
// Use ExtractExecution to convert its result into an Execution.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// Delete deletes the specified execution.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
