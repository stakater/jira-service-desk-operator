package controllers

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/types"

	"github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
)

var _ = Describe("Customer Controller", func() {

	ns, _ = os.LookupEnv("OPERATOR_NAMESPACE")

	AfterEach(func() {
		cUtil.TryDeleteCustomer(mockData.SampleCustomer.Spec.Name, ns)
		util.TryDeleteProject(mockData.CreateProjectInput.Spec.Name, ns)
	})

	Describe("Create new Jira Service Desk customer", func() {
		Context("With valid fields", func() {
			It("should create a new customer", func() {
				_ = cUtil.CreateCustomer(mockData.SampleCustomer, ns)
				customer := cUtil.GetCustomer(mockData.SampleCustomer.Spec.Name, ns)

				Expect(customer.Status.CustomerId).To(Equal(""))
			})
		})
	})

	// Describe("Add Jira Service Desk customer to project", func() {
	// 	Context("With Valid Project Id", func() {
	// 		It("Should add the customer in the project", func() {
	// 			//	_ = util.CreateProject(mockData.CreateProjectInput, ns)
	// 			project := util.GetProject(mockData.CreateProjectInput.Spec.Name, ns)

	// 			Expect(project.Status.ID).ToNot(Equal(""))

	// 			//_ = cUtil.CreateCustomer(mockData.SampleCustomer, ns)
	// 			customer := cUtil.GetCustomer(mockData.SampleCustomer.Spec.Name, ns)

	// 			Expect(customer.Status.CustomerId).ToNot(Equal(""))

	// 			customer.Spec.Projects = []string{"TEST"}

	// 			_ = cUtil.UpdateCustomer(customer, ns)
	// 			updatedCustomer := cUtil.GetCustomer(customer.Spec.Name, ns)

	// 			Expect(customer.Spec.Projects).To(Equal(updatedCustomer.Status.AssociatedProjects))
	// 		})
	// 	})
	// })

	// Describe("Remove Jira Service Desk customer from project", func() {
	// 	Context("With Valid Project Id", func() {
	// 		It("Should remove the customer from that project", func() {

	// 			//	_ = util.CreateProject(mockData.CreateProjectInput, ns)
	// 			project := util.GetProject(mockData.CreateProjectInput.Spec.Name, ns)

	// 			Expect(project.Status.ID).ToNot(Equal(""))

	// 			//				_ = cUtil.CreateCustomer(mockData.SampleCustomer, ns)
	// 			customer := cUtil.GetCustomer(mockData.SampleCustomer.Spec.Name, ns)

	// 			Expect(customer.Status.CustomerId).ToNot(Equal(""))

	// 			customer.Spec.Projects = []string{}

	// 			_ = cUtil.UpdateCustomer(customer, ns)
	// 			updatedCustomer := cUtil.GetCustomer(customer.Spec.Name, ns)

	// 			Expect(updatedCustomer.Status.AssociatedProjects).To(BeNil())
	// 		})
	// 	})
	// })

	Describe("Delete Jira Service Desk customer", func() {
		Context("With valid Customer AccountId", func() {
			It("should delete the customer", func() {
				_ = cUtil.CreateCustomer(mockData.SampleCustomer, ns)

				customer := cUtil.GetCustomer(mockData.SampleCustomer.Spec.Name, ns)
				Expect(customer.Status.CustomerId).NotTo(BeEmpty())

				cUtil.DeleteCustomer(customer.Name, ns)

				customerObject := &v1alpha1.Customer{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: mockData.SampleCustomer.Spec.Name, Namespace: ns}, customerObject)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
