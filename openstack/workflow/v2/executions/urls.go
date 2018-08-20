package executions

import (
	"github.com/gophercloud/gophercloud"
)

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("executions", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("executions")
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("executions")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("executions", id)
}
