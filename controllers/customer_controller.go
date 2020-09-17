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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	jiraservicedeskclient "github.com/stakater/jira-service-desk-operator/pkg/jiraservicedesk/client"
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
)

const (
	// TODO: Check if this is required in our case
	// defaultRequeueTime = 60 * time.Second
	CustomerFinalizer string = "jiraservicedesk.stakater.com/customer"
)

// CustomerReconciler reconciles a Customer object
type CustomerReconciler struct {
	client.Client
	Log                   logr.Logger
	Scheme                *runtime.Scheme
	JiraServiceDeskClient jiraservicedeskclient.Client
}

// +kubebuilder:rbac:groups=jiraservicedesk.stakater.com,resources=customers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jiraservicedesk.stakater.com,resources=customers/status,verbs=get;update;patch

func (r *CustomerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	log := r.Log.WithValues("customer", req.NamespacedName)

	log.Info("Reconciling Customer")

	// Fetch the Customer instance
	instance := &jiraservicedeskv1alpha1.Customer{}

	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading the object - donot requeue the request
		return reconcilerUtil.RequeueWithError(err)
	}

	// Validate Custom Resource
	if ok, err := instance.IsValid(); !ok {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	// Resource is marked for deletion
	if instance.DeletionTimestamp != nil {
		log.Info("Deletion timestamp found for instance " + req.Name)
		if finalizerUtil.HasFinalizer(instance, CustomerFinalizer) {
			return r.handleDelete(req, instance)
		}
		// Finalizer doesn't exist so clean up is already done
		return reconcilerUtil.DoNotRequeue()
	}

	// Add finalizer if it doesn't exist
	if !finalizerUtil.HasFinalizer(instance, CustomerFinalizer) {
		log.Info("Adding finalizer for instance" + req.Name)

		finalizerUtil.AddFinalizer(instance, CustomerFinalizer)

		err := r.Client.Update(context.TODO(), instance)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}

		return reconcilerUtil.DoNotRequeue()
	}

	return r.handleCreate(req, instance)
}

func (r *CustomerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jiraservicedeskv1alpha1.Customer{}).
		Complete(r)
}

func (r *CustomerReconciler) handleCreate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	if len(instance.Status.AccountId) == 0 {
		log.Info("Creating Jira Service Desk Customer: " + instance.Spec.DisplayName)

		customer := r.JiraServiceDeskClient.GetCustomerFromCustomerCRForCreateCustomer(instance)
		customerId, err := r.JiraServiceDeskClient.CreateCustomer(customer)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}

		instance.Status.AccountId = customerId
		log.Info("Successfully created Jira Service Desk Customer: " + instance.Spec.DisplayName)
	}

	for _, projectKey := range instance.Spec.ProjectKeys {
		if !contains(instance.Status.ProjectKeys, projectKey) {
			err := r.JiraServiceDeskClient.AddCustomerToProject(instance.Status.AccountId, projectKey)
			if err != nil {
				return reconcilerUtil.ManageError(r.Client, instance, err, false)
			}
			instance.Status.ProjectKeys = append(instance.Status.ProjectKeys, projectKey)
			log.Info("Successfully added Jira Service Desk Customer into project: " + projectKey)
		}
	}

	for index, projectKey := range instance.Status.ProjectKeys {
		if !contains(instance.Spec.ProjectKeys, projectKey) {
			err := r.JiraServiceDeskClient.RemoveCustomerFromProject(instance.Status.AccountId, projectKey)
			if err != nil {
				return reconcilerUtil.ManageError(r.Client, instance, err, false)
			}
			instance.Status.ProjectKeys[index] = ""
			log.Info("Successfully removed Jira Service Desk Customer from project: " + projectKey)
		}
	}

	instance.Status.ProjectKeys = deleteEmpty(instance.Status.ProjectKeys)

	return reconcilerUtil.ManageSuccess(r.Client, instance)
}

func (r *CustomerReconciler) handleDelete(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	if instance == nil {
		// Instance not found, nothing to do
		return reconcilerUtil.DoNotRequeue()
	}

	log.Info("Removing Jira Service Desk Customer: " + instance.Spec.DisplayName)

	// Remove customer from JSD project
	for _, projectKey := range instance.Status.ProjectKeys {
		err := r.JiraServiceDeskClient.RemoveCustomerFromProject(instance.Status.AccountId, projectKey)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		log.Info("Successfully removed Jira Service Desk Customer from project: " + projectKey)
	}

	// Delete Finalizer
	finalizerUtil.DeleteFinalizer(instance, CustomerFinalizer)

	log.Info("Finalizer removed for customer: " + instance.Spec.DisplayName)

	// Update instance
	err := r.Client.Update(context.TODO(), instance)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	return reconcilerUtil.DoNotRequeue()
}

func contains(slice []string, search string) bool {
	for _, value := range slice {
		if value == search {
			return true
		}
	}
	return false
}

func deleteEmpty(slice []string) []string {
	var output []string
	for _, str := range slice {
		if str != "" {
			output = append(output, str)
		}
	}
	return output
}
