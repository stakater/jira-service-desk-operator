package jiraservicedeskclient

import (
	"net/http"
)

const (
	EndpointApiVersion3Project = "/rest/api/3/project"
)

// TODO: Remove this code and populate these values from secrets
var (
	apiBaseUrl string
	apiToken   string
)

type Client interface {
}

// Client wraps http client
type jiraServiceDeskClient struct {
	apiToken   string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates an API client
func NewClient(apiToken string, baseURL string) Client {
	return &jiraServiceDeskClient{
		apiToken:   apiToken,
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}
