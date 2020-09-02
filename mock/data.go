package mock

import "strconv"

var DummyProjectIDStr = "1007"
var DummyProjectIDInt, _ = strconv.Atoi(DummyProjectIDStr)

var BaseURL = "https://mock.atlassian.net/"

var CreateProjectRequestJSON = map[string]string{
	"name":               "stakater",
	"key":                "STK",
	"projectTypeKey":     "service_desk",
	"projectTemplateKey": "com.atlassian.servicedesk:itil-v2-service-desk-project",
	"description":        "Sample project for jira-service-desk-operator",
	"assigneeType":       "PROJECT_LEAD",
	"leadAccountId":      "5ebfbc3wwe226gfda32c3590",
	"url":                "https://stakater.com",
}

var CreateProjectResponseJSON = map[string]interface{}{
	"self": BaseURL + "/rest/api/3/project/" + DummyProjectIDStr,
	"id":   DummyProjectIDInt,
	"key":  "STK",
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
	Name:               "stakater",
	Key:                "STK",
	ProjectTypeKey:     "service_desk",
	ProjectTemplateKey: "com.atlassian.servicedesk:itil-v2-service-desk-project",
	Description:        "Sample project for jira-service-desk-operator",
	AssigneeType:       "PROJECT_LEAD",
	LeadAccountId:      "5ebfbc3wwe226gfda32c3590",
	URL:                "https://stakater.com",
}

var UpdateProjectRequestJSON = map[string]string{
	"name": "stakater2",
	"key":  "WEE",
}

var UpdateProjectResponseJSON = map[string]interface{}{
	"self": BaseURL + "/rest/api/3/project/" + DummyProjectIDStr,
	"id":   DummyProjectIDInt,
	"key":  "STK",
}

var UpdateProjectInput = struct {
	Name string
	Key  string
}{
	Name: "stakater2",
	Key:  "WEE",
}
