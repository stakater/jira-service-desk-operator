package client

import (
	"errors"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"

	mockData "github.com/stakater/jira-service-desk-operator/mock"
)

func TestJiraService_GetProject_withValidId_shouldGetProjectByThatId(t *testing.T) {
	defer gock.Off()

	// The project.LeadAccountId & project.ProjectTemplateKey are not returned by the GET call. Therefore, those values are not returned and matched via Gock Server
	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Get("/10003").
		Reply(200).
		JSON(mockData.GetProjectByIdResponseJSON)

	jiraClient := NewClient("", mockData.BaseURL, "")
	project, err := jiraClient.GetProjectById("/10003")

	st.Expect(t, project.Description, mockData.GetProjectByIdExpectedResponse.Description)
	st.Expect(t, project.Name, mockData.GetProjectByIdExpectedResponse.Name)
	st.Expect(t, project.AssigneeType, mockData.GetProjectByIdExpectedResponse.AssigneeType)
	st.Expect(t, project.ProjectTypeKey, mockData.GetProjectByIdExpectedResponse.ProjectTypeKey)
	st.Expect(t, project.Key, mockData.GetProjectByIdExpectedResponse.Key)
	st.Expect(t, project.URL, mockData.GetProjectByIdExpectedResponse.URL)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_CreateProject_withValidData_shouldCreateProject(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Post("/").
		MatchType("json").
		JSON(mockData.CreateProjectInputJSON).
		Reply(200).
		JSON(map[string]interface{}{"self": "https://jira-service-desk.net/rest/api/3/project/10003", "id": 10003, "key": "KEY"})

	sampleProject := Project{
		Name:               mockData.CreateProjectInput.Name,
		Key:                mockData.CreateProjectInput.Key,
		ProjectTypeKey:     mockData.CreateProjectInput.ProjectTypeKey,
		ProjectTemplateKey: mockData.CreateProjectInput.ProjectTemplateKey,
		Description:        mockData.CreateProjectInput.Description,
		AssigneeType:       mockData.CreateProjectInput.AssigneeType,
		LeadAccountId:      mockData.CreateProjectInput.LeadAccountId,
		URL:                mockData.CreateProjectInput.URL,
	}

	jiraClient := NewClient("", mockData.BaseURL, "")
	id, err := jiraClient.CreateProject(sampleProject)

	st.Expect(t, id, "10003")
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_UpdateProject_withValidData_shouldUpdateProject(t *testing.T) {
}

func TestJiraService_DeleteProject_withValidId_shouldDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Delete("/10003").
		Reply(204)

	jiraClient := NewClient("", mockData.BaseURL, "")

	err := jiraClient.DeleteProject("10003")
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_DeleteProject_withInValidId_shouldNotDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointApiVersion3Project).
		Delete("/").
		Reply(404)

	jiraClient := NewClient("", mockData.BaseURL, "")

	err := jiraClient.DeleteProject("10003")

	st.Expect(t, err, errors.New("Rest request to delete Project failed with status: 404"))

	st.Expect(t, gock.IsDone(), true)
}
