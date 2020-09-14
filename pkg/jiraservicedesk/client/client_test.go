package client

import (
	"errors"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"

	"github.com/stakater/jira-service-desk-operator/mock"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
)

func TestJiraService_GetProject_shouldGetProject_whenValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Get("/" + mockData.ProjectID).
		Reply(200).
		JSON(mockData.GetProjectByIdResponseJSON)

	jiraClient := NewClient("", mockData.BaseURL, "")
	project, err := jiraClient.GetProjectById("/" + mockData.ProjectID)

	st.Expect(t, project.Description, mockData.GetProjectByIdExpectedResponse.Description)
	st.Expect(t, project.Name, mockData.GetProjectByIdExpectedResponse.Name)
	st.Expect(t, project.AssigneeType, mockData.GetProjectByIdExpectedResponse.AssigneeType)
	st.Expect(t, project.ProjectTypeKey, mockData.GetProjectByIdExpectedResponse.ProjectTypeKey)
	st.Expect(t, project.Key, mockData.GetProjectByIdExpectedResponse.Key)
	st.Expect(t, project.URL, mockData.GetProjectByIdExpectedResponse.URL)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_GetProject_shouldNotGetProject_whenInValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Get("/").
		Reply(404)

	jiraClient := NewClient("", mockData.BaseURL, "")
	_, err := jiraClient.GetProjectById("/" + mockData.ProjectID)

	st.Expect(t, err, errors.New(mockData.GetProjectFailedErrorMsg))
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_CreateProject_shouldCreateProject_whenValidProjectDataIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Post("/").
		MatchType("json").
		JSON(mockData.CreateProjectInputJSON).
		Reply(200).
		JSON(mockData.CreateProjectResponseJSON)

	sampleProject := Project{
		Name:               mockData.CreateProjectInput.Spec.Name,
		Key:                mockData.CreateProjectInput.Spec.Key,
		ProjectTypeKey:     mockData.CreateProjectInput.Spec.ProjectTypeKey,
		ProjectTemplateKey: mockData.CreateProjectInput.Spec.ProjectTemplateKey,
		Description:        mockData.CreateProjectInput.Spec.Description,
		AssigneeType:       mockData.CreateProjectInput.Spec.AssigneeType,
		LeadAccountId:      mockData.CreateProjectInput.Spec.LeadAccountId,
		URL:                mockData.CreateProjectInput.Spec.URL,
	}

	jiraClient := NewClient("", mockData.BaseURL, "")
	id, err := jiraClient.CreateProject(sampleProject)

	st.Expect(t, id, mockData.ProjectID)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_CreateProject_shouldNotCreateProject_whenInValidProjectDataIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Post("/").
		Reply(400)

	sampleProject := Project{
		Name: mockData.CreateProjectInput.Spec.Name,
	}

	jiraClient := NewClient("", mockData.BaseURL, "")
	_, err := jiraClient.CreateProject(sampleProject)

	st.Expect(t, err, errors.New(mockData.CreateProjectFailedErrorMsg))
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_UpdateProject_shouldUpdateProject_whenValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Put("/" + mockData.ProjectID).
		JSON(mock.UpdateProjectRequestJSON).
		Reply(204).
		JSON(mock.UpdateProjectResponseJSON)

	var updateProject = Project{
		Key:  mockData.UpdateProjectInput.Key,
		Name: mockData.UpdateProjectInput.Name,
	}

	client := NewClient("", mock.BaseURL, "")
	err := client.UpdateProject(updateProject, mock.ProjectID)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_UpdateProject_shouldNotUpdateProject_whenInvalidProjectDataIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mock.BaseURL + EndpointApiVersion3Project).
		Put("/").
		JSON(mock.UpdateProjectRequestJSON).
		Reply(404)
		//No Json is sent here, just checking the error

	var updateProject = Project{
		Key:  mock.UpdateProjectInput.Key,
		Name: mock.UpdateProjectInput.Name,
	}

	client := NewClient("", mock.BaseURL, "")
	err := client.UpdateProject(updateProject, mock.ProjectID)
	st.Expect(t, err, errors.New(mockData.UpdateProjectFailedErrorMsg))

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_DeleteProject_shouldDeleteProject_whenValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Delete("/" + mockData.ProjectID).
		Reply(204)

	jiraClient := NewClient("", mockData.BaseURL, "")

	err := jiraClient.DeleteProject(mockData.ProjectID)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_DeleteProject_shouldNotDeleteProject_whenInValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Delete("/").
		Reply(404)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.DeleteProject(mockData.ProjectID)

	st.Expect(t, err, errors.New(mockData.DeleteProjectFailedErrorMsg))
	st.Expect(t, gock.IsDone(), true)
}
