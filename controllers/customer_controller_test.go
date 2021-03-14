package controllers

import (
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/types"

	"github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
)

var _ = Describe("Customer Controller", func() {

	ns, _ = os.LookupEnv("OPERATOR_NAMESPACE")

	customerInput := mockData.SampleCustomer
	// Randomize customer name and email
	str := cUtil.RandSeqString(3)
	customerInput.Spec.Name += str
	customerInput.Spec.Email = "customer" + str + "@sample.com"
	customerInput.Spec.Projects = []string{strings.ToUpper(projectKey)}

	AfterEach(func() {
		cUtil.TryDeleteCustomer(customerInput.Spec.Name, ns)
	})

	Describe("Create new Jira Service Desk customer", func() {
		Context("With valid fields", func() {
			It("should create a new customer", func() {
				_ = cUtil.CreateCustomer(customerInput, ns)
				customer := cUtil.GetCustomer(customerInput.Spec.Name, ns)

				Expect(customer.Status.CustomerId).ToNot(Equal(""))
			})
		})
	})

	Describe("Modifying customer associations", func() {
		Describe("Add Jira Service Desk customer to project", func() {
			Context("With Valid Project Id", func() {
				It("Should add the customer in the project", func() {
					project := util.GetProject(mockData.CustomerTestProjectInput.Spec.Name, ns)
					Expect(project.Status.ID).ToNot(Equal(""))

					_ = cUtil.CreateCustomer(customerInput, ns)
					time.Sleep(5 * time.Second)

					customer := cUtil.GetCustomer(customerInput.Spec.Name, ns)

					Expect(customer.Status.CustomerId).ToNot(Equal(""))

					customer.Spec.Projects = []string{strings.ToUpper(customerKey)}

					_ = cUtil.UpdateCustomer(customer, ns)
					updatedCustomer := cUtil.GetCustomer(customer.Spec.Name, ns)

					Expect(customer.Spec.Projects).To(Equal(updatedCustomer.Status.AssociatedProjects))
				})
			})
		})

		Describe("Remove Jira Service Desk customer from project", func() {
			Context("With Valid Project Id", func() {
				It("Should remove the customer from that project", func() {
					project := util.GetProject(mockData.CustomerTestProjectInput.Spec.Name, ns)
					Expect(project.Status.ID).ToNot(Equal(""))

					mockData.SampleUpdatedCustomer.Spec.Name = customerInput.Spec.Name
					mockData.SampleUpdatedCustomer.Spec.Email = customerInput.Spec.Email
					// Assigning Customer -> CustomerTestproject Key
					mockData.SampleUpdatedCustomer.Spec.Projects = []string{strings.ToUpper(customerKey)}

					_ = cUtil.CreateCustomer(mockData.SampleUpdatedCustomer, ns)
					time.Sleep(5 * time.Second)

					customer := cUtil.GetCustomer(mockData.SampleUpdatedCustomer.Spec.Name, ns)

					Expect(customer.Status.CustomerId).ToNot(Equal(""))

					customer.Spec.Projects = []string{strings.ToUpper(projectKey)}

					_ = cUtil.UpdateCustomer(customer, ns)
					updatedCustomer := cUtil.GetCustomer(customer.Spec.Name, ns)

					Expect(customer.Spec.Projects).To(Equal(updatedCustomer.Status.AssociatedProjects))
				})
			})
		})
	})

	Describe("Delete Jira Service Desk customer", func() {
		Context("With valid Customer AccountId", func() {
			It("should delete the customer", func() {

				_ = cUtil.CreateCustomer(customerInput, ns)

				customer := cUtil.GetCustomer(customerInput.Spec.Name, ns)
				Expect(customer.Status.CustomerId).NotTo(BeEmpty())

				cUtil.DeleteCustomer(customer.Name, ns)

				customerObject := &v1alpha1.Customer{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: customerInput.Spec.Name, Namespace: ns}, customerObject)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
