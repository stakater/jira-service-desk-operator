package client

import jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"

func customerToCustomerCRMapper(customer Customer) jiraservicedeskv1alpha1.Customer {
	var customerObject jiraservicedeskv1alpha1.Customer

	customerObject.Spec.Email = customer.Email
	customerObject.Spec.Name = customer.DisplayName
	customerObject.Spec.Projects = customer.ProjectKeys

	return customerObject
}

func customerCRToCustomerMapperForCreateCustomer(customer *jiraservicedeskv1alpha1.Customer) Customer {
	customerObject := Customer{
		DisplayName: customer.Spec.Name,
		Email:       customer.Spec.Email,
	}
	return customerObject
}

func customerGetResponseToCustomerMapper(response CustomerGetResponse) Customer {
	return Customer{
		AccountId:   response.AccountId,
		DisplayName: response.DisplayName,
		Email:       response.EmailAddress,
	}
}
