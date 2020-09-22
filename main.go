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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	"github.com/stakater/jira-service-desk-operator/controllers"
	jiraservicedeskclient "github.com/stakater/jira-service-desk-operator/pkg/jiraservicedesk/client"
	jiraservicedeskconfig "github.com/stakater/jira-service-desk-operator/pkg/jiraservicedesk/config"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(jiraservicedeskv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	watchNamespace, err := getWatchNamespace()
	if err != nil {
		setupLog.Error(err, "unable to get WatchNamespace, "+
			"the manager will watch and manage resources in all Namespaces")
	}

	// Manager options
	options := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "48610e2b.stakater.com",
		Namespace:          watchNamespace, // namespaced-scope when the value is not an empty string
	}

	// Add support for MultiNamespace set in WATCH_NAMESPACE (e.g ns1,ns2)
	// More Info: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/cache#MultiNamespacedCacheBuilder
	if strings.Contains(watchNamespace, ",") {
		setupLog.Info("manager will be watching namespace %q", watchNamespace)
		// configure cluster-scoped with MultiNamespacedCacheBuilder
		options.Namespace = ""
		options.NewCache = cache.MultiNamespacedCacheBuilder(strings.Split(watchNamespace, ","))
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Load config for controller from secret
	controllerConfig, err := jiraservicedeskconfig.LoadControllerConfig(mgr.GetAPIReader())
	if err != nil {
		setupLog.Error(err, "unable to load controller config")
		os.Exit(1)
	}

	if err = (&controllers.ProjectReconciler{
		Client:                mgr.GetClient(),
		Scheme:                mgr.GetScheme(),
		Log:                   ctrl.Log.WithName("controllers").WithName("Project"),
		JiraServiceDeskClient: jiraservicedeskclient.NewClient(controllerConfig.ApiToken, controllerConfig.ApiBaseUrl, controllerConfig.Email),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Project")
		os.Exit(1)
	}

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&jiraservicedeskv1alpha1.Project{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Project")
			os.Exit(1)
		}
		if err = (&jiraservicedeskv1alpha1.Customer{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Customer")
			os.Exit(1)
		}
	}

	if err = (&controllers.CustomerReconciler{
		Client:                mgr.GetClient(),
		Log:                   ctrl.Log.WithName("controllers").WithName("Customer"),
		Scheme:                mgr.GetScheme(),
		JiraServiceDeskClient: jiraservicedeskclient.NewClient(controllerConfig.ApiToken, controllerConfig.ApiBaseUrl, controllerConfig.Email),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Customer")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// getWatchNamespace returns the Namespace the operator should be watching for changes
func getWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	var watchNamespaceEnvVar = "WATCH_NAMESPACE"

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}
