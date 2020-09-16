package client

import jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"

func customerCRToCustomerMapperForCreateCustomer(customer *jiraservicedeskv1alpha1.Customer) Customer {
	customerObject := Customer{
		DisplayName: customer.Spec.DisplayName,
		Email:       customer.Spec.Email,
	}
	return customerObject
}
