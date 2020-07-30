package client

import (
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	"net/http"
)

const (
	EndpointApiVersion3Project = "/rest/api/3/project"
)

type Client interface {
	// Methods for Project
	GetProjectByName(name string) (Project, error)
	GetProjectFromCR(spec jiraservicedeskv1alpha1.ProjectSpec) Project
	CreateProject(spec jiraservicedeskv1alpha1.ProjectSpec) (Project, error)
	UpdateProject(updatedProject Project) (Project, error)
	ProjectEqual(oldProject Project, newProject Project) bool
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
