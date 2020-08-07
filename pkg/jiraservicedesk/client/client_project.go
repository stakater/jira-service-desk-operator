package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

const (
	// Endpoints
	EndpointApiVersion3Project = "/rest/api/3/project"

	// Project Template Types
	ClassicProjectTemplateKey = "com.atlassian.servicedesk:itil-v2-service-desk-project"
	NextGenProjectTemplateKey = "com.atlassian.servicedesk:next-gen-it-service-desk"
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

type ProjectGetResponse struct {
	Self           string      `json:"self,omitempty"`
	Id             string      `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Key            string      `json:"key,omitempty"`
	Description    string      `json:"description,omitempty"`
	Lead           ProjectLead `json:"lead,omitempty"`
	ProjectTypeKey string      `json:"projectTypeKey,omitempty"`
	Style          string      `json:"style,omitempty"`
	AssigneeType   string      `json:"assigneeType,omitempty"`
	URL            string      `json:"url,omitempty"`
}

type ProjectLead struct {
	Self      string `json:"self,omitempty"`
	AccountId string `json:"accountId,omitempty"`
}

type ProjectCreateResponse struct {
	Self string `json:"self"`
	Id   int    `json:"id"`
	Key  string `json:"key"`
}

func (c *jiraServiceDeskClient) GetProjectById(id string) (Project, error) {
	var project Project

	request, err := c.newRequest("GET", EndpointApiVersion3Project+"/"+id, nil)
	if err != nil {
		return project, err
	}

	response, err := c.do(request)
	if err != nil {
		return project, err
	}
	defer response.Body.Close()

	var responseObject ProjectGetResponse
	err = json.NewDecoder(response.Body).Decode(&responseObject)
	if err != nil {
		return project, err
	}

	project = projectGetResponseToProjectMapper(responseObject)
	return project, err
}

func (c *jiraServiceDeskClient) CreateProject(project Project) (string, error) {
	request, err := c.newRequest("POST", EndpointApiVersion3Project, project)
	if err != nil {
		return "", err
	}

	response, err := c.do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	responseData, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err := errors.New("Rest request to create Project failed with status " + strconv.Itoa(response.StatusCode) +
			" and response: " + string(responseData))
		return "", err
	}

	var responseObject ProjectCreateResponse
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return "", err
	}
	projectId := strconv.Itoa(responseObject.Id)

	return projectId, err
}

func (c *jiraServiceDeskClient) UpdateProject(updatedProject Project) error {
	// Add logic for updating project here
	return nil
}

func (c *jiraServiceDeskClient) DeleteProject(id string) error {
	request, err := c.newRequest("DELETE", EndpointApiVersion3Project+"/"+id, nil)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 204 {
		return errors.New("Rest request to delete Project failed with status: " + strconv.Itoa(response.StatusCode))
	}

	return err
}

func (c *jiraServiceDeskClient) ProjectEqual(oldProject Project, newProject Project) bool {
	// The fields AvatarId, IssueSecurityScheme, NotificationScheme, PermissionScheme, CategoryId are not retrieved
	// through get project REST API call so they cannot be used in project comparison
	return oldProject.Id == newProject.Id &&
		oldProject.Name == newProject.Name &&
		oldProject.Key == newProject.Key &&
		oldProject.ProjectTypeKey == newProject.ProjectTypeKey &&
		oldProject.ProjectTemplateKey == newProject.ProjectTemplateKey &&
		oldProject.Description == newProject.Description &&
		oldProject.AssigneeType == newProject.AssigneeType &&
		oldProject.LeadAccountId == newProject.LeadAccountId &&
		oldProject.URL == newProject.URL
}

func (c *jiraServiceDeskClient) GetProjectFromProjectCR(project *jiraservicedeskv1alpha1.Project) Project {
	return projectCRToProjectMapper(project)
}
