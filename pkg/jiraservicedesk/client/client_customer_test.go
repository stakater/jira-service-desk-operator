package client

import (
	"testing"

	"github.com/nbio/st"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
	"gopkg.in/h2non/gock.v1"
)

func TestJiraClient_GetCustomerById_shouldGetCustomer_whenValidCustomerAccountIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL+EndpointUser).
		Get("/").
		MatchParam("accountId", mockData.CustomerAccountId).
		Reply(200).
		JSON(mockData.GetCustomerResponseJSON)

	jiraClient := NewClient("", mockData.BaseURL, "")
	customer, err := jiraClient.GetCustomerById(mockData.CustomerAccountId)

	st.Expect(t, customer.AccountId, mockData.GetCustomerResponse.AccountId)
	st.Expect(t, customer.DisplayName, mockData.GetCustomerResponse.DisplayName)
	st.Expect(t, customer.Email, mockData.GetCustomerResponse.Email)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_GetCustomerById_shouldNotGetCustomer_whenInValidCustomerAccountIdIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + EndpointUser).
		Get("/").
		Reply(400)

	jiraClient := NewClient("", mockData.BaseURL, "")
	customer, err := jiraClient.GetCustomerById(mockData.CustomerAccountId)

	st.Expect(t, customer.AccountId, "")
	st.Expect(t, customer.DisplayName, "")
	st.Expect(t, customer.Email, "")
	st.Reject(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_CreateCustomer_shouldCreateCustomer_whenValidCustomerDataIsGiven(t *testing.T) {
	defer gock.Off()

	gock.New(mockData.BaseURL + CreateCustomerApiPath).
		Post("/").
		MatchType("json").
		JSON(mockData.CreateCustomerInputJSON).
		Reply(201).
		JSON(mockData.CreateCustomerResponseJSON)

	sampleCustomer := Customer{
		Email:       mockData.GetCustomerResponse.Email,
		DisplayName: mockData.GetCustomerResponse.DisplayName,
	}

	jiraClient := NewClient("", mockData.BaseURL, "")
	id, err := jiraClient.CreateCustomer(sampleCustomer)

	st.Expect(t, id, mockData.CustomerAccountId)
	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_CreateCustomer_shouldNotCreateCustomer_whenInValidCustomerDataIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL + CreateCustomerApiPath).
		Post("/").
		Reply(400)

	sampleCustomer := Customer{
		DisplayName: mockData.GetCustomerResponse.DisplayName,
	}

	jiraClient := NewClient("", mockData.BaseURL, "")
	id, err := jiraClient.CreateCustomer(sampleCustomer)

	st.Expect(t, id, "")
	st.Reject(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_AddCustomerToProject_shouldAddCustomerToProject_whenValidProjectIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL + AddCustomerApiPath + mockData.AddProjectKey).
		Post(mockData.CustomerEndPoint).
		MatchType("json").
		JSON(map[string]interface{}{
			"accountIds": []string{mockData.CustomerAccountId},
		}).
		Reply(201)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.AddCustomerToProject(mockData.CustomerAccountId, mockData.AddProjectKey)

	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_AddCustomerToProject_shouldNotAddCustomerToProject_whenInValidProjectIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL + AddCustomerApiPath + mockData.AddProjectKey).
		Post(mockData.CustomerEndPoint).
		Reply(400)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.AddCustomerToProject(mockData.CustomerAccountId, mockData.AddProjectKey)

	st.Reject(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_RemoveCustomerFromProject_shouldRemoveCustomerFromProject_whenValidProjectIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL + AddCustomerApiPath + mockData.RemoveProjectKey).
		Post(mockData.CustomerEndPoint).
		Reply(201)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.AddCustomerToProject(mockData.CustomerAccountId, mockData.RemoveProjectKey)

	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_RemoveCustomerFromProject_shouldNotRemoveCustomerFromProject_whenInvalidProjectIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL + AddCustomerApiPath + mockData.RemoveProjectKey).
		Post(mockData.CustomerEndPoint).
		Reply(400)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.AddCustomerToProject(mockData.CustomerAccountId, mockData.RemoveProjectKey)

	st.Reject(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_DeleteCustomer_shouldDeleteCustomer_whenValidCustomerIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL+EndpointUser).
		Delete("/").
		MatchParam("accountId", mockData.CustomerAccountId).
		Reply(200)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.DeleteCustomer(mockData.CustomerAccountId)

	st.Expect(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestJiraClient_DeleteCustomer_shouldNotDeleteCustomer_whenInvalidCustomerIsGiven(t *testing.T) {
	defer gock.Off()
	gock.New(mockData.BaseURL + EndpointUser).
		Delete("/").
		Reply(400)

	jiraClient := NewClient("", mockData.BaseURL, "")
	err := jiraClient.DeleteCustomer(mockData.CustomerAccountId)

	st.Reject(t, err, nil)

	st.Expect(t, gock.IsDone(), true)
}
