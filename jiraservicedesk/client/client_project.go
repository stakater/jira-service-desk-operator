package client

import (
	"errors"
	"io/ioutil"
	"strconv"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

type Project struct {
	Id                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Key                 string `json:"key,omitempty"`
	ProjectTypeKey      string `json:"projectTypeKey,omitempty"`
	ProjectTemplateKey  string `json:"projectTemplateKey,omitempty"`
	Description         string `json:"description,omitempty"`
	AssigneeType        string `json:"assigneeType,omitempty"`
	LeadAccountId       string `json:"leadAccountId,omitempty"`
	URL                 string `json:"url,omitempty"`
	AvatarId            int    `json:"avatarId,omitempty"`
	IssueSecurityScheme int    `json:"issueSecurityScheme,omitempty"`
	PermissionScheme    int    `json:"permissionScheme,omitempty"`
	NotificationScheme  int    `json:"notificationScheme,omitempty"`
	CategoryId          int    `json:"categoryId,omitempty"`
}

func NewProject(name string) Project {
	return Project{
		Name: name,
	}
}

func (c *jiraServiceDeskClient) GetProjectByKey(name string) (Project, error) {
	return NewProject("test"), nil
}

func (c *jiraServiceDeskClient) CreateProject(project Project) error {
	request, err := c.newRequest("POST", EndpointApiVersion3Project, project)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		data, _ := ioutil.ReadAll(response.Body)
		err := errors.New("Rest request to create Project failed with status " + strconv.Itoa(response.StatusCode) +
			" and response: " + string(data))
		return err
	}

	return err
}

func (c *jiraServiceDeskClient) UpdateProject(updatedProject Project) (Project, error) {
	return NewProject("test"), nil
}

func (c *jiraServiceDeskClient) ProjectEqual(oldProject Project, newProject Project) bool {
	return false
}

func (c *jiraServiceDeskClient) GetProjectFromProjectSpec(spec jiraservicedeskv1alpha1.ProjectSpec) Project {
	return projectSpecToProjectMapper(spec)
}

func projectSpecToProjectMapper(spec jiraservicedeskv1alpha1.ProjectSpec) Project {
	return Project{
		Name:                spec.Name,
		Key:                 spec.Key,
		ProjectTypeKey:      spec.ProjectTypeKey,
		ProjectTemplateKey:  spec.ProjectTemplateKey,
		Description:         spec.Description,
		AssigneeType:        spec.AssigneeType,
		LeadAccountId:       spec.LeadAccountId,
		URL:                 spec.URL,
		AvatarId:            spec.AvatarId,
		IssueSecurityScheme: spec.IssueSecurityScheme,
		PermissionScheme:    spec.IssueSecurityScheme,
		NotificationScheme:  spec.NotificationScheme,
		CategoryId:          spec.CategoryId,
	}
}
