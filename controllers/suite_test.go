/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	controllerUtil "github.com/stakater/jira-service-desk-operator/controllers/util"
	mockData "github.com/stakater/jira-service-desk-operator/mock"
	c "github.com/stakater/jira-service-desk-operator/pkg/jiraservicedesk/client"
	"github.com/stakater/jira-service-desk-operator/pkg/jiraservicedesk/config"
	secretsUtil "github.com/stakater/operator-utils/util/secrets"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

var ctx context.Context
var r *ProjectReconciler
var util *controllerUtil.TestUtil
var ns string

var cr *CustomerReconciler
var cUtil *controllerUtil.TestUtil

var log = logf.Log.WithName("config")

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {

	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = jiraservicedeskv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})

	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	ctx = context.Background()

	// Retrieve operator namespace
	ns, _ := os.LookupEnv("OPERATOR_NAMESPACE")

	apiToken, err := secretsUtil.LoadSecretDataUsingClient(k8sClient, config.JiraServiceDeskSecretName, ns, config.JiraServiceDeskAPITokenSecretKey)
	Expect(err).ToNot(HaveOccurred())
	Expect(apiToken).ToNot(BeNil())

	apiBaseUrl, err := secretsUtil.LoadSecretDataUsingClient(k8sClient, config.JiraServiceDeskSecretName, ns, config.JiraServiceDeskAPIBaseURLSecretKey)
	Expect(err).ToNot(HaveOccurred())
	Expect(apiBaseUrl).ToNot(BeNil())

	email, err := secretsUtil.LoadSecretDataUsingClient(k8sClient, config.JiraServiceDeskSecretName, ns, config.JiraServiceDeskEmailSecretKey)
	Expect(err).ToNot(HaveOccurred())
	Expect(email).ToNot(BeNil())

	r = &ProjectReconciler{
		Client:                k8sClient,
		Scheme:                scheme.Scheme,
		Log:                   log.WithName("Reconciler"),
		JiraServiceDeskClient: c.NewClient(apiToken, apiBaseUrl, email),
	}
	Expect(r).ToNot((BeNil()))

	util = controllerUtil.New(ctx, k8sClient, r)
	Expect(util).ToNot(BeNil())

	cr = &CustomerReconciler{
		Client:                k8sClient,
		Scheme:                scheme.Scheme,
		Log:                   log.WithName("Reconciler"),
		JiraServiceDeskClient: c.NewClient(apiToken, apiBaseUrl, email),
	}
	Expect(cr).ToNot((BeNil()))

	cUtil = controllerUtil.New(ctx, k8sClient, cr)
	Expect(util).ToNot(BeNil())

	_ = util.CreateProject(mockData.CustomerTestProjectInput, ns)

	close(done)
}, 60)

var _ = AfterSuite(func() {
	util.TryDeleteProject(mockData.CustomerTestProjectInput.Spec.Name, ns)

	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
