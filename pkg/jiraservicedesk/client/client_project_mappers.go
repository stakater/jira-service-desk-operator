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
		PermissionScheme:    project.Spec.PermissionScheme,
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

func projectToProjectCRMapper(project Project) jiraservicedeskv1alpha1.Project {

	var projectObject jiraservicedeskv1alpha1.Project

	projectObject.Status.ID = project.Id
	projectObject.Spec.Name = project.Name
	projectObject.Spec.Key = project.Key
	projectObject.Spec.ProjectTypeKey = project.ProjectTypeKey
	projectObject.Spec.ProjectTemplateKey = project.ProjectTemplateKey
	projectObject.Spec.Description = project.Description
	projectObject.Spec.AssigneeType = project.AssigneeType
	projectObject.Spec.LeadAccountId = project.LeadAccountId
	projectObject.Spec.URL = project.URL
	projectObject.Spec.AvatarId = project.AvatarId
	projectObject.Spec.IssueSecurityScheme = project.IssueSecurityScheme
	projectObject.Spec.PermissionScheme = project.PermissionScheme
	projectObject.Spec.NotificationScheme = project.NotificationScheme
	projectObject.Spec.CategoryId = project.CategoryId

	return projectObject
}
