package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

const (
	// Endpoints
	EndpointCreateCustomer = "/rest/servicedeskapi/customer"
	EndpointAddCustomer    = "/rest/servicedeskapi/servicedesk/"
)

type Customer struct {
	AccountId   string   `json:"accountId,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Email       string   `json:"email,omitempty"`
	ProjectKeys []string `json:"projectKeys,omitempty"`
}

type CustomerCreateResponse struct {
	AccountId    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

type AddCustomerResponse struct {
	AccountIds []string `json:"accountIds,omitempty"`
}

func (c *jiraServiceDeskClient) CreateCustomer(customer Customer) (string, error) {
	request, err := c.newRequest("POST", EndpointCreateCustomer, customer)
	if err != nil {
		return "", err
	}

	response, err := c.do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	responseData, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = errors.New("Rest request to create Customer failed with status " + strconv.Itoa(response.StatusCode) +
			" and response: " + string(responseData))
		return "", err
	}

	var responseObject CustomerCreateResponse
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return "", err
	}

	return responseObject.AccountId, err
}

func (c *jiraServiceDeskClient) AddCustomerToProject(customerAccountId string, projectKey string) error {
	addCustomerBody := AddCustomerResponse{
		AccountIds: []string{customerAccountId},
	}

	request, err := c.newRequest("POST", EndpointAddCustomer+projectKey+"/customer", addCustomerBody)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = errors.New("Rest request to add Customer failed with status " + strconv.Itoa(response.StatusCode))
		return err
	}

	return nil
}

func (c *jiraServiceDeskClient) RemoveCustomerFromProject(customerAccountId string, projectKey string) error {
	removeCustomerBody := AddCustomerResponse{
		AccountIds: []string{customerAccountId},
	}

	request, err := c.newRequestWithExperimentalHeader("DELETE", EndpointAddCustomer+projectKey+"/customer", removeCustomerBody)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = errors.New("Rest request to remove Customer failed with status " + strconv.Itoa(response.StatusCode))
		return err
	}

	return nil
}

func (c *jiraServiceDeskClient) GetCustomerFromCustomerCRForCreateCustomer(customer *jiraservicedeskv1alpha1.Customer) Customer {
	return customerCRToCustomerMapperForCreateCustomer(customer)
}
