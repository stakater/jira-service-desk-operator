package client

// func TestJiraClient_GetCustomerById_shouldGetCustomer_whenValidCustomerAccountIdIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	gock.New(mockData.BaseURL + "/rest/api/3/user?accountId=").
// 		Get("sample12345").
// 		Reply(200).
// 		JSON(map[string]string{
// 			"self":         "https://sample.net/user?accountId=sample12345",
// 			"accountId":    "sample12345",
// 			"emailAddress": "sample@test.com",
// 			"displayName":  "Sample Customer",
// 			"accountType":  "customer",
// 		})

// 	jiraClient := NewClient("", mockData.BaseURL, "")
// 	customer, err := jiraClient.GetCustomerById("sample12345")

// 	st.Expect(t, customer.AccountId, "sample12345")
// 	st.Expect(t, customer.DisplayName, "Sample Customer")
// 	st.Expect(t, customer.Email, "sample@test.com")
// 	st.Expect(t, err, nil)

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_GetCustomerById_shouldNotGetCustomer_whenInValidCustomerAccountIdIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	gock.New(mockData.BaseURL)

// 	st.Expect(t, gock.IsDone(), true)
// }

// var CustomerAccountId = "sample12345"

// var CreateCustomerInput = Customer{
// 	DisplayName: "Sample Customer",
// 	Email:       "sample@test.com",
// }

// var CreateCustomerInputJSON = map[string]string{
// 	"displayName": "Sample Customer",
// 	"email":       "sample@test.com",
// }

// var CreateCustomerResponseJSON = map[string]string{
// 	"accountId":    "sample12345",
// 	"displayName":  "Sample Customer",
// 	"emailAddress": "sample@test.com",
// }

// func TestJiraClient_CreateCustomer_shouldCreateCustomer_whenValidCustomerDataIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	gock.New(mockData.BaseURL + CreateCustomerApiPath).
// 		Post("/").
// 		MatchType("json").
// 		JSON(CreateCustomerInputJSON).
// 		Reply(201).
// 		JSON(CreateCustomerResponseJSON)

// 	jiraClient := NewClient("", mockData.BaseURL, "")
// 	id, err := jiraClient.CreateCustomer(CreateCustomerInput)

// 	st.Expect(t, id, CustomerAccountId)
// 	st.Expect(t, err, nil)

// 	// Verify that we don't have pending mocks
// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_CreateCustomer_shouldNotCreateCustomer_whenInValidCustomerDataIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_AddCustomerToProject_shouldAddCustomerToProject_whenValidProjectIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_AddCustomerToProject_shouldNotAddCustomerToProject_whenInValidProjectIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_RemoveCustomerFromProject_shouldRemoveCustomerFromProject_whenValidProjectIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_RemoveCustomerFromProject_shouldNotRemoveCustomerFromProject_whenInvalidProjectIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_DeleteCustomer_shouldDeleteCustomer_whenValidCustomerIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }

// func TestJiraClient_DeleteCustomer_shouldNotDeleteCustomer_whenInvalidCustomerIsGiven(t *testing.T) {
// 	defer gock.Off()

// 	st.Expect(t, gock.IsDone(), true)
// }
