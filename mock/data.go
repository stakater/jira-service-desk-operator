package mock

import (
	"strconv"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

const BaseURL = "https://sample.atlassian.net"

var ProjectID = "10003"
var ProjectIDInt, _ = strconv.Atoi(ProjectID)
var InvalidPermissionScheme = "4000"

var GetProjectFailedErrorMsg = "Rest request to get Project failed with status: 404"
var CreateProjectFailedErrorMsg = "Rest request to create Project failed with status 400 and response: "
var UpdateProjectFailedErrorMsg = "Rest request to update Project failed with status 404 and response: "
var DeleteProjectFailedErrorMsg = "Rest request to delete Project failed with status: 404"

var CreateProjectInputJSON = map[string]string{
	"name":               "testproject",
	"key":                "TEST",
	"projectTypeKey":     "service_desk",
	"projectTemplateKey": "com.atlassian.servicedesk:itil-v2-service-desk-project",
	"description":        "Sample project for jira-service-desk-operator",
	"assigneeType":       "PROJECT_LEAD",
	"leadAccountId":      "5f62e5902b42470070d1fb83",
	"url":                "https://test.com",
}

var CreateProjectResponseJSON = map[string]interface{}{
	"self": BaseURL + "/rest/api/3/project/10003",
	"id":   ProjectIDInt,
	"key":  "KEY",
}

var CreateProjectInput = jiraservicedeskv1alpha1.Project{
	Spec: jiraservicedeskv1alpha1.ProjectSpec{
		Name:               "testproject",
		Key:                "TEST",
		ProjectTypeKey:     "service_desk",
		ProjectTemplateKey: "com.atlassian.servicedesk:itil-v2-service-desk-project",
		Description:        "Sample project for jira-service-desk-operator",
		AssigneeType:       "PROJECT_LEAD",
		LeadAccountId:      "5f62e5902b42470070d1fb83",
		URL:                "https://test.com",
	},
}

var CreateProjectInvalidInput = jiraservicedeskv1alpha1.Project{
	Spec: jiraservicedeskv1alpha1.ProjectSpec{
		Name:                "test",
		Key:                 "TEST20000",
		ProjectTypeKey:      "service_desk",
		ProjectTemplateKey:  "com.atlassian.servicedesk:itil-v2-service-desk-project",
		Description:         "Sample project for jira-service-desk-operator",
		AssigneeType:        "PROJECT_LEAD",
		LeadAccountId:       "5f62e5902b42470070d1fb83",
		URL:                 "https://test.com",
		AvatarId:            10200,
		IssueSecurityScheme: 10001,
		PermissionScheme:    10011,
		NotificationScheme:  10021,
		CategoryId:          10120,
	},
}

var UpdateMutableProjectFields = struct {
	Name string
	Key  string
}{
	"testupdated",
	"TEST2",
}

var UpdateImmutableProjectFields = struct {
	ProjectTypeKey string
}{
	"business",
}

var GetProjectByIdResponseJSON = map[string]string{
	"description":    "Sample Project",
	"name":           "Sample",
	"assigneeType":   "UNASSIGNED",
	"projectTypeKey": "business",
	"key":            "KEY",
	"url":            "https://www.sample.com",
}

var GetProjectByIdExpectedResponse = struct {
	Description    string
	Name           string
	AssigneeType   string
	ProjectTypeKey string
	Key            string
	URL            string
}{
	"Sample Project",
	"Sample",
	"UNASSIGNED",
	"business",
	"KEY",
	"https://www.sample.com",
}

var UpdateProjectInput = struct {
	Id   string
	Name string
	Key  string
}{
	Id:   "99999",
	Name: "stakater2",
	Key:  "WEE",
}

var UpdateProjectRequestJSON = map[string]string{
	"name": "stakater2",
	"key":  "WEE",
}

var UpdateProjectResponseJSON = map[string]interface{}{
	"self": BaseURL + "/rest/api/3/project/" + ProjectID,
	"id":   ProjectIDInt,
	"key":  "STK",
}
