package client

import (
	"errors"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

var endpoint = "177"
var client = NewClient("", "https://mock.atlassian.net/", "")
var jira_url = "https://mock.atlassian.net/rest/api/3"

var project = Project{
	Name:               "Stakater",
	Key:                "STK",
	ProjectTypeKey:     "service_desk",
	ProjectTemplateKey: "com.atlassian.servicedesk:itil-v2-service-desk-project",
	Description:        "Sample project for jira-service-desk-operator",
	AssigneeType:       "PROJECT_LEAD",
	LeadAccountId:      "5ebfbc3wwe226gfda32c3590",
	URL:                "https://stakater.com",
}

func TestJiraServiceDesk_DeleteProject_withValidProjectId_shouldDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New(jira_url).
		Delete("/" + endpoint).
		Reply(204)

	err := client.DeleteProject(endpoint)
	st.Expect(t, err, nil)

	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_DeleteProject_withInvaildProjectId_shouldNotDeleteProject(t *testing.T) {
	defer gock.Off()

	gock.New(jira_url).
		Delete("/" + endpoint).
		Reply(210)

	err := client.DeleteProject(endpoint + "3")

	st.Expect(t, err, errors.New("Rest request to delete Project failed with status: 210"))

	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_CreateProject__shouldCreateProject(t *testing.T) {
	defer gock.Off()

	gock.New(jira_url).
		Post("").
		JSON(map[string]string{"name": "Stakater",
			"key":                "STK",
			"projectTypeKey":     "service_desk",
			"projectTemplateKey": "com.atlassian.servicedesk:itil-v2-service-desk-project",
			"description":        "Sample project for jira-service-desk-operator",
			"assigneeType":       "PROJECT_LEAD",
			"leadAccountId":      "5ebfbc3wwe226gfda32c3590",
			"url":                "https://stakater.com"}).
		Reply(204).
		JSON(map[string]interface{}{"self": "https://mock.atlassian.net/rest/api/3/project/1007",
			"id":  1007,
			"key": "STK"})

	projectId, err := client.CreateProject(project)
	st.Expect(t, err, nil)
	st.Expect(t, projectId, "1007")
	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}
