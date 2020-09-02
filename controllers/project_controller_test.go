package controllers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	mockData "github.com/stakater/jira-service-desk-operator/mock"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("ProjectController", func() {

	var name string
	var key string
	var projectTypeKey string
	var projectTemplateKey string
	var description string
	var assigneeType string
	var leadAccountId string
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

	AfterEach(func() {
		util.TryDeleteProject(name, ns)
	})

	Describe("Create New JiraServiceDeskProject Resource", func() {
		Context("With the required fields", func() {
			It("should create a new project", func() {
				_ = util.CreateProject(name, key, projectTypeKey, projectTemplateKey, description, assigneeType, leadAccountId, url, ns)
				project := util.GetProject(name, ns)

				Expect(project.Status.ID).NotTo(BeEmpty())
			})
		})
	})

	Describe("Deleting jira service desk project resource", func() {
		Context("When project on jira service desk was created", func() {
			It("should remove resource and delete project ", func() {
				_ = util.CreateProject(name, key, projectTypeKey, projectTemplateKey, description, assigneeType, leadAccountId, url, ns)
				project := util.GetProject(name, ns)
				Expect(project.Status.ID).NotTo(BeEmpty())
				util.DeleteProject(name, ns)
				projectObject := &jiraservicedeskv1alpha1.Project{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, projectObject)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Updating jira service desk resource", func() {
		Context("With mutable fields ", func() {
			It("should assign changed field values to Project", func() {
				_ = util.CreateProject(name, key, projectTypeKey, projectTemplateKey, description, assigneeType, leadAccountId, url, ns)
				project := util.GetProject(name, ns)
				project.Spec.Name = mockData.UpdateMutableProjectFields.Name
				project.Spec.Key = mockData.UpdateMutableProjectFields.Key
				err := k8sClient.Update(ctx, project)
				if err != nil {
					Fail(err.Error())
				}
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ns}}
				_, err = r.Reconcile(req)
				if err != nil {
					Fail(err.Error())
				}
				updatedProject := util.GetProject(name, ns)
				Expect(updatedProject.Spec.Name).To(Equal(mockData.UpdateMutableProjectFields.Name))
				Expect(updatedProject.Spec.Key).To(Equal(mockData.UpdateMutableProjectFields.Key))
			})
		})
		Context("With immutable fields ", func() {
			It("should not assign changed field values to Project", func() {
				_ = util.CreateProject(name, key, projectTypeKey, projectTemplateKey, description, assigneeType, leadAccountId, url, ns)
				project := util.GetProject(name, ns)
				project.Spec.ProjectTypeKey = mockData.UpdateImmutableProjectFields.ProjectTypeKey
				err := k8sClient.Update(ctx, project)
				if err != nil {
					Fail(err.Error())
				}
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: ns}}
				_, err = r.Reconcile(req)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
