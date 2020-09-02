package util

import (
	"context"

	"github.com/onsi/ginkgo"
	ginko "github.com/onsi/ginkgo"
	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

//CreateJiraProjectObject creates a jira service desk custom resource object
func (t *TestUtil) CreateJiraProjectObject(name string, key string, projectTypeKey string, projectTemplateKey string, description string, assigneeType string, leadAccountId string, url string, namespace string) *jiraservicedeskv1alpha1.Project {
	return &jiraservicedeskv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: jiraservicedeskv1alpha1.ProjectSpec{
			Name:               name,
			Key:                key,
			ProjectTypeKey:     projectTypeKey,
			ProjectTemplateKey: projectTemplateKey,
			Description:        description,
			AssigneeType:       assigneeType,
			LeadAccountId:      leadAccountId,
			URL:                url,
		},
	}
}

// CreateProject creates and submits a jira service desk project object to the kubernetes server
func (t *TestUtil) CreateProject(name string, key string, projectTypeKey string, projectTemplateKey string, description string, assigneeType string, leadAccountId string, url string, namespace string) *jiraservicedeskv1alpha1.Project {
	projectObject := t.CreateJiraProjectObject(name, key, projectTypeKey, projectTemplateKey, description, assigneeType, leadAccountId, url, namespace)
	err := t.k8sClient.Create(t.ctx, projectObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginko.Fail(err.Error())
	}

	return projectObject
}

// GetProject fetches a project object from kubernetes
func (t *TestUtil) GetProject(name string, namespace string) *jiraservicedeskv1alpha1.Project {
	channelObject := &jiraservicedeskv1alpha1.Project{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, channelObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	return channelObject
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

// TryDeleteProject - Tries to delete Project if it exists, does not fail on any error
func (t *TestUtil) TryDeleteProject(name string, namespace string) {
	projectObject := &jiraservicedeskv1alpha1.Project{}
	_ = t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)
	_ = t.k8sClient.Delete(t.ctx, projectObject)
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, _ = t.r.Reconcile(req)
}
