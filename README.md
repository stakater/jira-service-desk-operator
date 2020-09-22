# jira-service-desk-operator
Kubernetes operator for Jira Service Desk

## About

Jira service desk(JSD) operator is used to automate the process of setting up JSD for alertmanager in a k8s native way. By using CRDs it lets you:

1. Manage Projects
2. Manage customer/organization for projects
3. Configure Issues

It uses [Jira REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/) in it's underlying layer and can be extended to perform other tasks that are supported via the REST API.

## Usage

### Prerequisites

- Atlassian account
- API Token to access Jira REST API (https://id.atlassian.com/manage-profile/security/api-tokens)

### Create secret

Create the following secret which is required for jira-service-desk-operator:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: jira-service-desk-config
  namespace: default
data:
  JIRA_SERVICE_DESK_API_TOKEN: <API_TOKEN>
  #Example: https://stakater-cloud.atlassian.net/
  JIRA_SERVICE_DESK_API_BASE_URL: <JSD_BASE_URL>
  JIRA_SERVICE_DESK_EMAIL: <EMAIL>
type: Opaque
```

### Deploy operator

- Make sure that [certman](https://cert-manager.io/) is deployed in your cluster since webhooks require certman to generate valid certs since webhooks serve using HTTPS
- To install certman
```terminal
$ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.16.1/cert-manager.yaml
```
- Deploy operator
```terminal
$ oc apply -f bundle/manifests
```

### Project
We support the following CRUD operation on project via our Jira Service Desk Operator:
* Create - Creates a new projects with the provided fields
* Update - Updates an existing project with the updated fields
* Delete - Removes and deletes the project 

Examples for Project Custom Resource can be found at [here](https://github.com/stakater/jira-service-desk-operator/tree/master/examples/project).

#### Limitations:
* We only support creating three types of JSD projects via our operator i.e Business, ServiceDesk, Software. The details and differences between these project types can be viewed [here](https://confluence.atlassian.com/adminjiraserver/jira-applications-and-project-types-overview-938846805.html).
* Following are the immutable fields that cannot be updated:
    * ProjectTemplateKey
    * ProjectTypeKey
    * leadAccountId 
    * CategoryId 
    * NotificationScheme
    * PermissionScheme 
    * issueSecurityScheme 

    You can read more about these fields on [Jira Service Desk api docs](https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-projects/#api-rest-api-3-project-post).


### Customer:
We support the following CRUD operations on customer via our Jira Service Desk Operator
* Create - Create a new customer and assign the projects mentioned in the CR
* Update - Only updates(add/remove) the associated projects mentioned in the CR
* Delete - Remove all the project associations and deletes the customer

Examples for Project Custom Resource can be found at [here](https://github.com/stakater/jira-service-desk-operator/blob/handle-customers/examples/customer/customer.yaml).

#### Limitations:
* Jira Service Desk Operator can access only those customers which are created through it. Customers that are manually created and added in the projects can’t be accessed later with the Jira Service Desk Operator.
* Each custom resource is associated to a single customer. 
* You can not update **customer name and email**.

## Local Development

[Operator-sdk v0.19.0](https://github.com/operator-framework/operator-sdk/releases/tag/v0.19.0) is required for local development.

1. Create `jira-service-desk-config` secret
2. Run `make run ENABLE_WEBHOOKS=false WATCH_NAMESPACE=default OPERATOR_NAMESPACE=default` where `WATCH_NAMESPACE` denotes the namespaces that the operator is supposed to watch and `OPERATOR_NAMESPACE` is the namespace in which it's supposed to be deployed.

3. Before committing your changes run the following to ensure that everything is verified and up-to-date:
   - `make verify`
   - `make bundle`
   - `make packagemanifests`
   
## Running Tests

### Pre-requisites:
1. Create a namespace with the name `test`
2. Create `jira-service-desk-config` secret in test namespace

### To run tests:
Use the following command to run tests:
`make test OPERATOR_NAMESPACE=test USE_EXISTING_CLUSTER=true`


