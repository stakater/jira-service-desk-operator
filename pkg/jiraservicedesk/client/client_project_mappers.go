package client

import (
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

func projectSpecToProjectMapper(spec jiraservicedeskv1alpha1.ProjectSpec) Project {
	return Project{
		Name:                spec.Name,
		Key:                 spec.Key,
		ProjectTypeKey:      spec.ProjectTypeKey,
		ProjectTemplateKey:  spec.ProjectTemplateKey,
		Description:         spec.Description,
		AssigneeType:        spec.AssigneeType,
		LeadAccountId:       spec.LeadAccountId,
		URL:                 spec.URL,
		AvatarId:            spec.AvatarId,
		IssueSecurityScheme: spec.IssueSecurityScheme,
		PermissionScheme:    spec.IssueSecurityScheme,
		NotificationScheme:  spec.NotificationScheme,
		CategoryId:          spec.CategoryId,
	}
}

func projectGetResponseToProjectMapper(response ProjectGetResponse) Project {
	var projectTemplateKey string
	if len(response.Style) > 0 {
		if response.Style == "classic" {
			projectTemplateKey = ClassicProjectTemplateKey
		} else if response.Style == "next-gen" {
			projectTemplateKey = NextGenProjectTemplateKey
		}
	}

	return Project{
		Id:                 response.Id,
		Name:               response.Name,
		Key:                response.Key,
		ProjectTypeKey:     response.ProjectTypeKey,
		ProjectTemplateKey: projectTemplateKey,
		Description:        response.Description,
		AssigneeType:       response.AssigneeType,
		LeadAccountId:      response.Lead.AccountId,
		URL:                response.URL,
	}
}
