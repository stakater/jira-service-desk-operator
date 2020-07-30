package client

import (
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

type Project struct {
	Name string `json:"name"`
}

func NewProject(name string) Project {
	return Project{
		Name: name,
	}
}

func (c *jiraServiceDeskClient) GetProjectByName(name string) (Project, error) {
	return NewProject("test"), nil
}

func (c *jiraServiceDeskClient) CreateProject(project Project) (Project, error) {
	return NewProject("test"), nil
}

func (c *jiraServiceDeskClient) UpdateProject(updatedProject Project) (Project, error) {
	return NewProject("test"), nil
}

func (c *jiraServiceDeskClient) ProjectEqual(oldProject Project, newProject Project) bool {
	return false
}

func (c *jiraServiceDeskClient) GetProjectFromSpec(spec jiraservicedeskv1alpha1.ProjectSpec) Project {
	return NewProject("test")
}
