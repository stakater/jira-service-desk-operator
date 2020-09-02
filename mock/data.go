package mock

import "strconv"

const BaseURL = "https://sample.atlassian.net"

var CreateProjectResponseID = "10003"
var CreateProjectResponseIDInt, _ = strconv.Atoi(CreateProjectResponseID)

var CreateProjectInputJSON = map[string]string{
	"name":               "testproject",
	"key":                "TEST",
	"projectTypeKey":     "service_desk",
	"projectTemplateKey": "com.atlassian.servicedesk:itil-v2-service-desk-project",
	"description":        "Sample project for jira-service-desk-operator",
	"assigneeType":       "PROJECT_LEAD",
	"leadAccountId":      "5ebfbc3ead226b0ba46c3590",
	"url":                "https://test.com",
}

var CreateProjectResponseJSON = map[string]interface{}{
	"self": BaseURL + "/rest/api/3/project/10003",
	"id":   CreateProjectResponseIDInt,
	"key":  "KEY",
}

var CreateProjectInput = struct {
	Name               string
	Key                string
	ProjectTypeKey     string
	ProjectTemplateKey string
	Description        string
	AssigneeType       string
	LeadAccountId      string
	URL                string
}{
	Name:               "testproject",
	Key:                "TEST",
	ProjectTypeKey:     "service_desk",
	ProjectTemplateKey: "com.atlassian.servicedesk:itil-v2-service-desk-project",
	Description:        "Sample project for jira-service-desk-operator",
	AssigneeType:       "PROJECT_LEAD",
	LeadAccountId:      "5ebfbc3ead226b0ba46c3590",
	URL:                "https://test.com",
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
