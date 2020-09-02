package client

import (
	"errors"
	"testing"

	"github.com/nbio/st"
	mock "github.com/stakater/jira-service-desk-operator/mock"
	"gopkg.in/h2non/gock.v1"
)

func TestJiraServiceDesk_DeleteProject_shouldDeleteProject_whenValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mock.BaseURL + EndpointApiVersion3Project).
		Delete("/" + mock.DummyProjectIDStr).
		Reply(204)

	client := NewClient("", mock.BaseURL, "")
	err := client.DeleteProject(mock.DummyProjectIDStr)
	st.Expect(t, err, nil)

	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_DeleteProject_shouldNotDeleteProject_whenInvaildProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mock.BaseURL + EndpointApiVersion3Project).
		Delete("/").
		Reply(404)
	client := NewClient("", mock.BaseURL, "")
	err := client.DeleteProject(mock.DummyProjectIDStr)

	st.Expect(t, err, errors.New("Rest request to delete Project failed with status: 404"))

	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_CreateProject_shouldCreateProject(t *testing.T) {
	defer gock.Off()

	gock.New(mock.BaseURL + EndpointApiVersion3Project).
		Post("").
		JSON(mock.CreateProjectRequestJSON).
		Reply(204).
		JSON(mock.CreateProjectResponseJSON)

	var project = Project{
		Name:               mock.CreateProjectInput.Name,
		Key:                mock.CreateProjectInput.Key,
		ProjectTypeKey:     mock.CreateProjectInput.ProjectTypeKey,
		ProjectTemplateKey: mock.CreateProjectInput.ProjectTemplateKey,
		Description:        mock.CreateProjectInput.Description,
		AssigneeType:       mock.CreateProjectInput.AssigneeType,
		LeadAccountId:      mock.CreateProjectInput.LeadAccountId,
		URL:                mock.CreateProjectInput.URL,
	}

	client := NewClient("", mock.BaseURL, "")
	projectID, err := client.CreateProject(project)
	st.Expect(t, err, nil)
	st.Expect(t, projectID, mock.DummyProjectIDStr)
	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_UpdateProject_shouldUpdateProject_whenValidProjectIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mock.BaseURL + EndpointApiVersion3Project).
		Put("/" + mock.DummyProjectIDStr).
		JSON(mock.UpdateProjectRequestJSON).
		Reply(204).
		JSON(mock.UpdateProjectResponseJSON)

	var updateProject = Project{
		Key:  mock.UpdateProjectInput.Key,
		Name: mock.UpdateProjectInput.Name,
	}
	client := NewClient("", mock.BaseURL, "")
	err := client.UpdateProject(updateProject, mock.DummyProjectIDStr)
	st.Expect(t, err, nil)
	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_UpdateProject_shouldNotUpdateProject_whenInvalidProjectIdIsGiven(t *testing.T) {
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
	err := client.UpdateProject(updateProject, mock.DummyProjectIDStr)
	st.Expect(t, err, errors.New("Rest request to update Project failed with status 404 and response: "))
	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}
