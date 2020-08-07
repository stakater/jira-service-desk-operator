package client

import (
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

func projectCRToProjectMapper(project *jiraservicedeskv1alpha1.Project) Project {

	projectObject := Project{
		Name:                project.Spec.Name,
		Key:                 project.Spec.Key,
		ProjectTypeKey:      project.Spec.ProjectTypeKey,
		ProjectTemplateKey:  project.Spec.ProjectTemplateKey,
		Description:         project.Spec.Description,
		AssigneeType:        project.Spec.AssigneeType,
		LeadAccountId:       project.Spec.LeadAccountId,
		URL:                 project.Spec.URL,
		AvatarId:            project.Spec.AvatarId,
		IssueSecurityScheme: project.Spec.IssueSecurityScheme,
		PermissionScheme:    project.Spec.IssueSecurityScheme,
		NotificationScheme:  project.Spec.NotificationScheme,
		CategoryId:          project.Spec.CategoryId,
	}

	if len(project.Status.ID) > 0 {
		projectObject.Id = project.Status.ID
	}

	return projectObject
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
