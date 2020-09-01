package controllers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
)

var _ = Describe("ProjectController", func() {

	Describe("Create New JiraServiceDeskProject Resource", func() {

		var description string
		var leadAccountId string
		var projectTemplateKey string
		var name string
		var assigneeType string
		var projectTypeKey string
		var key string
		var url string

		BeforeEach(func() {
			name = mockData.CreateProjectInput.Name
			key = mockData.CreateProjectInput.Key
			projectTypeKey = mockData.CreateProjectInput.ProjectTypeKey
			projectTemplateKey = mockData.CreateProjectInput.ProjectTemplateKey
			description = mockData.CreateProjectInput.Description
			assigneeType = mockData.CreateProjectInput.AssigneeType
			leadAccountId = mockData.CreateProjectInput.LeadAccountId
			url = mockData.CreateProjectInput.URL
		})

		Context("With the required fields", func() {
			It("should create a new project", func() {
				_ = util.CreateProject(name, key, projectTypeKey, projectTemplateKey, description, assigneeType, leadAccountId, url, ns)
				project := util.GetProject(name, ns)

				Expect(project.Status.ID).To(Equal(mockData.CreateProjectResponseID))
			})
		})
	})

})
