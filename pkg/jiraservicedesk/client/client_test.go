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

var update_project = Project{
	Key:  "WEE",
	Name: "Stakater2",
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
		Delete("/").
		Reply(404)

	err := client.DeleteProject("1999")

	st.Expect(t, err, errors.New("Rest request to delete Project failed with status: 404"))

	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

func TestJiraServiceDesk_CreateProject_shouldCreateProject(t *testing.T) {
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

// not done yet
func TestJiraServiceDesk_CreateProject_withInvalidProjectField_shouldNotCreateProject(t *testing.T) {
	defer gock.Off()

	gock.New(jira_url).
		Post("").
		JSON(map[string]string{"name": "Stakater",
			"keen":               "STK", // instead of key sending keen
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

func TestJiraServiceDesk_UpdateProject_withValidProjectId_shouldUpdateProject(t *testing.T) {
	defer gock.Off()

	gock.New(jira_url).
		Put("/1007").
		JSON(map[string]string{"name": "Stakater2",
			"key": "WEE"}).
		Reply(204).
		JSON(map[string]interface{}{"self": "https://mock.atlassian.net/rest/api/3/project/1007",
			"id":  1007,
			"key": "STK"})

	err := client.UpdateProject(update_project, "1007")
	st.Expect(t, err, nil)
	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}

//needs changes
func TestJiraServiceDesk_UpdateProject_withInvalidProjectId_shouldNotUpdateProject(t *testing.T) {
	defer gock.Off()

	gock.New(jira_url).
		Put("/").
		JSON(map[string]string{"name": "Stakater2",
			"key": "WEE"}).
		Reply(404)
		//JSON(map[string]interface{}{})

	err := client.UpdateProject(update_project, "1007")
	st.Expect(t, err, errors.New("Rest request to update Project failed with status 404 and response: "))
	// Verify no mock pending requests
	st.Expect(t, gock.IsDone(), true)
}
