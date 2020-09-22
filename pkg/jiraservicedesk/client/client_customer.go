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
	CreateCustomerApiPath = "/rest/servicedeskapi/customer"
	AddCustomerApiPath    = "/rest/servicedeskapi/servicedesk/"
	EndpointUser          = "/rest/api/3/user?accountId="
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

type CustomerAddResponse struct {
	AccountIds []string `json:"accountIds,omitempty"`
}

type CustomerGetResponse struct {
	Self         string `json:"self,omitempty"`
	AccountId    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	AccountType  string `json:"accountType,omitempty"`
}

func (c *jiraServiceDeskClient) GetCustomerById(customerAccountId string) (Customer, error) {
	var customer Customer

	request, err := c.newRequest("GET", EndpointUser+customerAccountId, nil, false)
	if err != nil {
		return customer, err
	}

	response, err := c.do(request)
	if err != nil {
		return customer, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err := errors.New("Rest request to get customer failed with status: " + strconv.Itoa(response.StatusCode))
		return customer, err
	}

	var responseObject CustomerGetResponse
	err = json.NewDecoder(response.Body).Decode(&responseObject)
	if err != nil {
		return customer, err
	}

	customer = customerGetResponseToCustomerMapper(responseObject)

	return customer, err
}

func (c *jiraServiceDeskClient) CreateCustomer(customer Customer) (string, error) {
	request, err := c.newRequest("POST", CreateCustomerApiPath, customer, false)
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
	addCustomerBody := CustomerAddResponse{
		AccountIds: []string{customerAccountId},
	}

	request, err := c.newRequest("POST", AddCustomerApiPath+projectKey+"/customer", addCustomerBody, false)
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
	removeCustomerBody := CustomerAddResponse{
		AccountIds: []string{customerAccountId},
	}

	request, err := c.newRequest("DELETE", AddCustomerApiPath+projectKey+"/customer", removeCustomerBody, true)
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

func (c *jiraServiceDeskClient) DeleteCustomer(customerAccountId string) error {
	request, err := c.newRequest("DELETE", EndpointUser+customerAccountId, nil, false)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = errors.New("Rest request to delete Customer failed with status " + strconv.Itoa(response.StatusCode))
		return err
	}

	return nil
}

func (c *jiraServiceDeskClient) GetCustomerCRFromCustomer(customer Customer) jiraservicedeskv1alpha1.Customer {
	return customerToCustomerCRMapper(customer)
}

func (c *jiraServiceDeskClient) GetCustomerFromCustomerCRForCreateCustomer(customer *jiraservicedeskv1alpha1.Customer) Customer {
	return customerCRToCustomerMapperForCreateCustomer(customer)
}
