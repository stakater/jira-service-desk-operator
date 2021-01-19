package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var Log = logf.Log.WithName("jiraServiceDeskClient")

type Client interface {
	// Methods for Project
	GetProjectById(id string) (Project, error)
	GetProjectFromProjectCR(project *jiraservicedeskv1alpha1.Project) Project
	GetProjectCRFromProject(project Project) jiraservicedeskv1alpha1.Project
	CreateProject(project Project) (string, error)
	DeleteProject(id string) error
	UpdateProject(updatedProject Project, id string) error
	ProjectEqual(oldProject Project, newProject Project) bool
	GetProjectForUpdateRequest(existingProject Project, newProject *jiraservicedeskv1alpha1.Project) Project
	UpdateProjectAccessPermissions(status bool, key string) error
	GetCustomerById(customerAccountId string) (Customer, error)
	CreateCustomer(customer Customer) (string, error)
	CreateLegacyCustomer(email string, projectKey string) (string, error)
	IsCustomerUpdated(customer *jiraservicedeskv1alpha1.Customer, existingCustomer Customer) bool
	AddCustomerToProject(customerAccountId string, projectKey string) error
	RemoveCustomerFromProject(customerAccountId string, projectKey string) error
	DeleteCustomer(customerAccountId string) error
	GetCustomerCRFromCustomer(customer Customer) jiraservicedeskv1alpha1.Customer
	GetCustomerFromCustomerCRForCreateCustomer(customer *jiraservicedeskv1alpha1.Customer) Customer
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

func (c *jiraServiceDeskClient) newRequest(method, path string, body interface{}, experimental bool) (*http.Request, error) {
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

	if experimental {
		req.Header.Set("X-ExperimentalApi", "opt-in")
	}

	req.SetBasicAuth(c.email, c.apiToken)
	return req, nil
}

func (c *jiraServiceDeskClient) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("Error calling the API endpoint: %v", err)
	}

	return resp, nil
}
