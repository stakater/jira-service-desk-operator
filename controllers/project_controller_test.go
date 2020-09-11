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

	projectInput := mockData.CreateProjectInput

	AfterEach(func() {
		util.TryDeleteProject(projectInput.Spec.Name, ns)
	})

	Describe("Create New JiraServiceDeskProject Resource", func() {
		Context("With the required fields", func() {
			It("should create a new project", func() {
				_ = util.CreateProject(projectInput, ns)
				project := util.GetProject(projectInput.Spec.Name, ns)

				Expect(project.Status.ID).NotTo(BeEmpty())
			})
		})
	})

	Describe("Deleting jira service desk project resource", func() {
		Context("When project on jira service desk was created", func() {
			It("should remove resource and delete project ", func() {
				_ = util.CreateProject(projectInput, ns)

				project := util.GetProject(projectInput.Spec.Name, ns)
				Expect(project.Status.ID).NotTo(BeEmpty())

				util.DeleteProject(project.Name, ns)

				projectObject := &jiraservicedeskv1alpha1.Project{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: projectInput.Spec.Name, Namespace: ns}, projectObject)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Updating jira service desk resource", func() {
		Context("With mutable fields ", func() {
			It("should assign changed field values to Project", func() {
				_ = util.CreateProject(projectInput, ns)
				project := util.GetProject(projectInput.Spec.Name, ns)

				project.Spec.Name = mockData.UpdateMutableProjectFields.Name
				project.Spec.Key = mockData.UpdateMutableProjectFields.Key

				err := k8sClient.Update(ctx, project)
				if err != nil {
					Fail(err.Error())
				}

				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: projectInput.Spec.Name, Namespace: ns}}
				_, err = r.Reconcile(req)
				if err != nil {
					Fail(err.Error())
				}

				updatedProject := util.GetProject(projectInput.Spec.Name, ns)

				Expect(updatedProject.Spec.Name).To(Equal(mockData.UpdateMutableProjectFields.Name))
				Expect(updatedProject.Spec.Key).To(Equal(mockData.UpdateMutableProjectFields.Key))
			})
		})
	})
})
