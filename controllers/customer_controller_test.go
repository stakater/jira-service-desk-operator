package controllers

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Customer Controller", func() {

	ns, _ = os.LookupEnv("OPERATOR_NAMESPACE")

	sampleCustomer := v1alpha1.Customer{
		Spec: v1alpha1.CustomerSpec{
			Name:  "sample",
			Email: "sample@email.com",
		},
	}

	AfterEach(func() {
		cUtil.TryDeleteCustomer(sampleCustomer.Spec.Name, ns)
	})

	Describe("Create new Jira Service Desk customer", func() {
		Context("With valid fields", func() {
			It("should create a new customer", func() {
				_ = cUtil.CreateCustomer(sampleCustomer, ns)
				customer := cUtil.GetCustomer(sampleCustomer.Spec.Name, ns)

				Expect(customer.Status.CustomerId).ToNot(Equal(""))
			})
		})
	})

	Describe("Delete Jira Service Desk customer", func() {
		Context("With valid Customer AccountId", func() {
			It("should delete the customer", func() {
				_ = cUtil.CreateCustomer(sampleCustomer, ns)

				customer := cUtil.GetCustomer(sampleCustomer.Spec.Name, ns)
				Expect(customer.Status.CustomerId).NotTo(BeEmpty())

				cUtil.DeleteCustomer(customer.Name, ns)

				customerObject := &v1alpha1.Customer{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: sampleCustomer.Spec.Name, Namespace: ns}, customerObject)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
