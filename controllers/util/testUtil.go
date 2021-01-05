package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/onsi/ginkgo"
	ginko "github.com/onsi/ginkgo"
	mockdata "github.com/stakater/jira-service-desk-operator/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

// TestUtil contains necessary objects required to perform operations during tests
type TestUtil struct {
	ctx       context.Context
	k8sClient client.Client
	r         reconcile.Reconciler
}

// New creates new TestUtil
func New(ctx context.Context, k8sClient client.Client, r reconcile.Reconciler) *TestUtil {
	return &TestUtil{
		ctx:       ctx,
		k8sClient: k8sClient,
		r:         r,
	}
}

// RandSeqString Generates a letter sequence with `n` characters
func (t *TestUtil) RandSeqString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// CreateNamespace creates a namespace in the kubernetes server
func (t *TestUtil) CreateNamespace(name string) {
	namespaceObject := t.CreateNamespaceObject(name)
	err := t.k8sClient.Create(t.ctx, namespaceObject)

	if err != nil {
		ginkgo.Fail(err.Error())
	}
}

// CreateNamespaceObject creates a namespace object
func (t *TestUtil) CreateNamespaceObject(name string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// DeleteNamespace deletes a namespace
func (t *TestUtil) DeleteNamespace(name string) {
	namespaceObject := &v1.Namespace{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name}, namespaceObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	err = t.k8sClient.Delete(t.ctx, namespaceObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
}

// CreateProjectObject creates a jira project custom resource object
func (t *TestUtil) CreateProjectObject(project jiraservicedeskv1alpha1.Project, namespace string) *jiraservicedeskv1alpha1.Project {
	return &jiraservicedeskv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      project.Spec.Name,
			Namespace: namespace,
		},
		Spec: jiraservicedeskv1alpha1.ProjectSpec{
			Name:                project.Spec.Name,
			Key:                 project.Spec.Key,
			ProjectTypeKey:      project.Spec.ProjectTypeKey,
			ProjectTemplateKey:  project.Spec.ProjectTemplateKey,
			Description:         project.Spec.Description,
			AssigneeType:        project.Spec.AssigneeType,
			LeadAccountId:       project.Spec.LeadAccountId,
			URL:                 project.Spec.URL,
			AvatarId:            project.Spec.AvatarId,
			IssueSecurityScheme: project.Spec.IssueSecurityScheme,
			PermissionScheme:    project.Spec.PermissionScheme,
			NotificationScheme:  project.Spec.NotificationScheme,
			CategoryId:          project.Spec.CategoryId,
		},
	}
}

// CreateCustomerObject creates a jira customer customer resource object
func (t *TestUtil) CreateCustomerObject(customer jiraservicedeskv1alpha1.Customer, namespace string) *jiraservicedeskv1alpha1.Customer {
	return &jiraservicedeskv1alpha1.Customer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      customer.Spec.Name,
			Namespace: namespace,
		},
		Spec: jiraservicedeskv1alpha1.CustomerSpec{
			Name:     customer.Spec.Name,
			Email:    customer.Spec.Email,
			Projects: customer.Spec.Projects,
		},
	}
}

// CreateProject creates and submits a Project object to the kubernetes server
func (t *TestUtil) CreateProject(project jiraservicedeskv1alpha1.Project, namespace string) *jiraservicedeskv1alpha1.Project {
	projectObject := t.CreateProjectObject(project, namespace)

	err := t.k8sClient.Create(t.ctx, projectObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: project.Spec.Name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return projectObject
}

// CreateCustomer creates and submits a new Customer object to the kubernetes server
func (t *TestUtil) CreateCustomer(customer jiraservicedeskv1alpha1.Customer, namespace string) *jiraservicedeskv1alpha1.Customer {
	customerObject := t.CreateCustomerObject(customer, namespace)

	err := t.k8sClient.Create(t.ctx, customerObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: customer.Spec.Name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return customerObject
}

// UpdateCustomer submits an updatedCustomer to the kubernetes server
func (t *TestUtil) UpdateCustomer(customerObject *jiraservicedeskv1alpha1.Customer, namespace string) *jiraservicedeskv1alpha1.Customer {

	err := t.k8sClient.Update(t.ctx, customerObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: customerObject.Spec.Name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return customerObject
}

// GetProject fetches a project object from kubernetes
func (t *TestUtil) GetProject(name string, namespace string) *jiraservicedeskv1alpha1.Project {
	projectObject := &jiraservicedeskv1alpha1.Project{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	return projectObject
}

// GetCustomer fetches a customer object from kubernetes
func (t *TestUtil) GetCustomer(name string, namespace string) *jiraservicedeskv1alpha1.Customer {
	customerObject := &jiraservicedeskv1alpha1.Customer{}

	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, customerObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return customerObject
}

// DeleteProject deletes the project resource
func (t *TestUtil) DeleteProject(name string, namespace string) {
	projectObject := &jiraservicedeskv1alpha1.Project{}

	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)
	if err != nil {
		ginko.Fail(err.Error())
	}

	err = t.k8sClient.Delete(t.ctx, projectObject)
	if err != nil {
		ginko.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, err = t.r.Reconcile(req)
	if err != nil {
		ginko.Fail(err.Error())
	}
}

// DeleteCustomer deletes the customer resource
func (t *TestUtil) DeleteCustomer(name string, namespace string) {
	customerObject := &jiraservicedeskv1alpha1.Customer{}

	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, customerObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	err = t.k8sClient.Delete(t.ctx, customerObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, err = t.r.Reconcile(req)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
}

// TryDeleteProject - Tries to delete Project if it exists, does not fail on any error
func (t *TestUtil) TryDeleteProject(name string, namespace string) {
	projectObject := &jiraservicedeskv1alpha1.Project{}
	_ = t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)
	_ = t.k8sClient.Delete(t.ctx, projectObject)
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, _ = t.r.Reconcile(req)
}

// TryDeleteCustomer - Tries to delete Customer if it exists, does not fail on any error
func (t *TestUtil) TryDeleteCustomer(name string, namespace string) {
	customerObject := &jiraservicedeskv1alpha1.Customer{}
	_ = t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, customerObject)
	_ = t.k8sClient.Delete(t.ctx, customerObject)
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, _ = t.r.Reconcile(req)
}

// DeleteAllProjects delete all the projects in the namespace
func (t *TestUtil) DeleteAllProjects(namespace string) {
	// Specify namespace in list Options
	listOptions := &client.ListOptions{Namespace: namespace}

	// List projects in a specified namespace
	projectList := &jiraservicedeskv1alpha1.ProjectList{}
	err := t.k8sClient.List(context.TODO(), projectList, listOptions)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	for _, project := range projectList.Items {
		project.Finalizers = []string{}

		err := t.k8sClient.Update(t.ctx, &project)
		if err != nil {
			if err.Error() == fmt.Sprintf(mockdata.ProjectObjectModifiedError, project.Name) {
				currentProject := t.GetProject(project.Name, namespace)
				currentProject.Finalizers = []string{}
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			} else {
				ginkgo.Fail(err.Error())
			}
		}

		t.TryDeleteProject(project.Name, namespace)
	}
}

// DeleteAllCustomers delete all the customers in the namespace
func (t *TestUtil) DeleteAllCustomers(namespace string) {
	// Specify namespace in list Options
	listOptions := &client.ListOptions{Namespace: namespace}

	// List customers in a specified namespace
	customerList := &jiraservicedeskv1alpha1.CustomerList{}
	err := t.k8sClient.List(context.TODO(), customerList, listOptions)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	for _, customer := range customerList.Items {
		customer.Finalizers = []string{}

		err := t.k8sClient.Update(t.ctx, &customer)
		if err != nil {
			if err.Error() == fmt.Sprintf(mockdata.CustomerObjectModifiedError, customer.Name) {
				currentCustomer := t.GetCustomer(customer.Name, namespace)
				currentCustomer.Finalizers = []string{}
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			} else {
				ginkgo.Fail(err.Error())
			}
		}

		t.TryDeleteCustomer(customer.Name, namespace)
	}
}
