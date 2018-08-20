package crontriggers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extension to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToCronTriggerListQuery() (string, error)
}

// ListOpts filters the result returned by the List() function.
type ListOpts struct {
	// WorkflowName allows to filter by workflow name.
	WorkflowName string `q:"workflow_name"`
	// WorkflowID allows to filter by workflow id.
	WorkflowID string `q:"workflow_id"`

	// Name allows to filter by trigger name.
	Name string `q:"name"`

	// Scope filters by the trigger's scope.
	// Values can be "private" or "public".
	Scope string `q:"scope"`

	// SortDir allows to select sort direction.
	// It can be "asc" or "desc" (default).
	SortDir string `q:"sort_dir"`

	// SortKey allows to sort by one of the cron trigger attributes.
	SortKey string `q:"sort_key"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of
	// cron triggers.
	Limit int `q:"limit"`
}

// ToCronTriggerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToCronTriggerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List performs a call to list cron triggers.
// You may provide options to filter the results.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToCronTriggerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return CronTriggerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extension to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToCronTriggerCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters used to create a cron trigger.
type CreateOpts struct {
	// Name is the cron trigger name.
	Name string `json:"name"`

	// Pattern is a Unix crontab patterns format to execute the workflow.
	Pattern string `json:"pattern"`

	// RemainingExecutions sets the number of executions for the trigger.
	RemainingExecutions int `json:"remaining_executions,omitempty"`

	// WorkflowInput is a JSON structure containing workflow input values, serialized as string.
	WorkflowInput string `json:"workflow_input,omitempty"`

	// WorkflowID is the unique id of the workflow.
	WorkflowID string `json:"workflow_id"`

	// WorkflowParams define workflow type specific parameters.
	WorkflowParams string `json:"workflow_params,omitempty"`
}

// ToCronTriggerCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToCronTriggerCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "execution")
}

// Create requests the creation of a new execution.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCronTriggerCreateMap()
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
// Use ExtractCronTrigger to convert its result into an CronTrigger.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// Delete deletes the specified execution.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
