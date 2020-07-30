package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	"io"
	"net/http"
	"net/url"
)

const (
	EndpointApiVersion3Project = "/rest/api/3/project"
)

type Client interface {
	// Methods for Project
	GetProjectByKey(key string) (Project, error)
	GetProjectFromProjectSpec(spec jiraservicedeskv1alpha1.ProjectSpec) Project
	CreateProject(project Project) (Project, error)
	UpdateProject(updatedProject Project) (Project, error)
	ProjectEqual(oldProject Project, newProject Project) bool
}

// Client wraps http client
type jiraServiceDeskClient struct {
	apiToken   string
	baseURL    string
	email      string
	httpClient *http.Client
}

// NewClient creates an API client
func NewClient(apiToken string, baseURL string, email string) Client {
	return &jiraServiceDeskClient{
		apiToken:   apiToken,
		baseURL:    baseURL,
		email:      email,
		httpClient: http.DefaultClient,
	}
}

func (c *jiraServiceDeskClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
	endpoint := c.baseURL + path
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "golang httpClient")
	req.SetBasicAuth(c.apiToken, "")
	return req, nil
}

func (c *jiraServiceDeskClient) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("Error calling the API endpoint: %v", err)
	}

	return resp, nil
}
