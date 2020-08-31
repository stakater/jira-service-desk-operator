package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

const BaseURL = "https://sample.atlassian.net/"

func TestJiraService_GetProject_shouldGetProjectById(t *testing.T) {
	defer gock.Off()

	// The project.LeadAccountId & project.ProjectTemplateKey are not returned by the GET call. Therefore, those values are not returned and matched via Gock Server
	gock.New(BaseURL + "rest/api/3/project").
		Get("/10003").
		Reply(200).
		JSON(map[string]string{
			"description":    "Sample Project",
			"name":           "Sample",
			"assigneeType":   "UNASSIGNED",
			"projectTypeKey": "business",
			"key":            "KEY",
			"url":            "https://www.sample.com",
		})

	jiraClient := NewClient("", BaseURL, "")
	project, err := jiraClient.GetProjectById("/10003")

	st.Expect(t, project.Description, "Sample Project")
	st.Expect(t, project.Name, "Sample")
	st.Expect(t, project.AssigneeType, "UNASSIGNED")
	st.Expect(t, project.ProjectTypeKey, "business")
	st.Expect(t, project.Key, "KEY")
	st.Expect(t, project.URL, "https://www.sample.com")
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_CreateProject_shouldCreateProject(t *testing.T) {
	defer gock.Off()

	gock.New(BaseURL + "/rest/api/3/project").
		Post("/").
		MatchType("json").
		JSON(map[string]string{
			"description":        "Sample Project",
			"leadAccountId":      "5ebfbc3ead226b0ba46c3591",
			"projectTemplateKey": "com.atlassian.jira.jira-incident-management-plugin:im-incident-management",
			"name":               "Sample",
			"assigneeType":       "UNASSIGNED",
			"projectTypeKey":     "business",
			"key":                "KEY",
		}).
		Reply(200).
		JSON(map[string]interface{}{"self": "https://jira-service-desk.net/rest/api/3/project/10003", "id": 10003, "key": "KEY"})

	sampleProject := Project{
		Description:        "Sample Project",
		LeadAccountId:      "5ebfbc3ead226b0ba46c3591",
		ProjectTemplateKey: "com.atlassian.jira.jira-incident-management-plugin:im-incident-management",
		Name:               "Sample",
		AssigneeType:       "UNASSIGNED",
		ProjectTypeKey:     "business",
		Key:                "KEY",
	}

	jiraClient := NewClient("", BaseURL, "")
	id, err := jiraClient.CreateProject(sampleProject)

	fmt.Println(err)

	st.Expect(t, id, "10003")
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_UpdateProject_shouldUpdateProject(t *testing.T) {
}

func TestJiraService_DeleteProject_shouldDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New(BaseURL + "/rest/api/3/project").
		Delete("/10003").
		Reply(204)

	jiraClient := NewClient("", BaseURL, "")

	err := jiraClient.DeleteProject("10003")
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraService_DeleteProject_shouldNotDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New(BaseURL + "/rest/api/3/project").
		Delete("/").
		Reply(404)

	jiraClient := NewClient("", BaseURL, "")

	err := jiraClient.DeleteProject("10003")

	st.Expect(t, err, errors.New("Rest request to delete Project failed with status: 404"))

	st.Expect(t, gock.IsDone(), true)
}
