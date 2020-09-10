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
)

type Customer struct {
	AccountId   string `json:"accountId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	//	projectKeys []string `json:"projectKeys,omitempty"`
}

type CustomerCreateResponse struct {
	AccountId    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
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

	if response.StatusCode != 201 {
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

func (c *jiraServiceDeskClient) GetCustomerFromCustomerCR(customer *jiraservicedeskv1alpha1.Customer) Customer {
	return customerCRToCustomerMapper(customer)
}
