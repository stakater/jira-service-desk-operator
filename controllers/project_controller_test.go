package controllers

import (
	"context"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Project Controller", func() {

	ns, _ = os.LookupEnv("OPERATOR_NAMESPACE")

	Describe("Positive test cases", func() {

		projectInput := mockData.CreateProjectInput

		// Generation of 3 char long random string
		key := cUtil.RandSeqString(3)

		projectInput.Spec.Name += key
		projectInput.Spec.Key = strings.ToUpper(key)

		AfterEach(func() {
			util.TryDeleteProject(projectInput.Spec.Name, ns)
		})

		Describe("Create new Jira service desk project resource", func() {
			Context("With valid fields", func() {
				It("should create a new project", func() {
					_ = util.CreateProject(projectInput, ns)
					project := util.GetProject(projectInput.Spec.Name, ns)

					Expect(project.Status.ID).ToNot(Equal(""))
				})
			})
		})

		Describe("Deleting jira service desk project resource", func() {
			Context("With valid project Id", func() {
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
					_, err = r.Reconcile(context.Background(), req)
					if err != nil {
						Fail(err.Error())
					}

					updatedProject := util.GetProject(projectInput.Spec.Name, ns)

					Expect(updatedProject.Spec.Name).To(Equal(mockData.UpdateMutableProjectFields.Name))
					Expect(updatedProject.Spec.Key).To(Equal(mockData.UpdateMutableProjectFields.Key))
				})

			})
			Context("With immutable fields ", func() {

				It("should not assign changed field values to Project", func() {
					_ = util.CreateProject(projectInput, ns)
					project := util.GetProject(projectInput.Spec.Name, ns)

					oldTypeKey := project.Spec.ProjectTypeKey
					project.Spec.ProjectTemplateKey = mockData.UpdateImmutableProjectFields.ProjectTypeKey

					err := k8sClient.Update(ctx, project)
					if err != nil {
						Fail(err.Error())
					}

					req := reconcile.Request{NamespacedName: types.NamespacedName{Name: projectInput.Spec.Name, Namespace: ns}}
					_, err = r.Reconcile(context.Background(), req)
					if err != nil {
						Fail(err.Error())
					}

					updatedProject := util.GetProject(projectInput.Spec.Name, ns)

					Expect(updatedProject.Spec.ProjectTypeKey).To(Equal(oldTypeKey))
				})
			})
		})
	})

	Describe("Negative test cases", func() {

		projectInvalidInput := mockData.CreateProjectInvalidInput

		Describe("Create new Jira servie desk project resource", func() {
			Context("with invalid fields", func() {
				It("should not create a new project", func() {
					key := cUtil.RandSeqString(9)
					projectInvalidInput.Spec.Key = strings.ToUpper(key)
					projectInvalidInput.Spec.Name += key[:3]

					_ = util.CreateProject(projectInvalidInput, ns)
					project := util.GetProject(projectInvalidInput.Spec.Name, ns)

					Expect(project.Status.ID).To(Equal(""))

					util.TryDeleteProject(projectInvalidInput.Spec.Name, ns)
				})
			})
		})
	})
})
