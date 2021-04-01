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
	"strings"

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
	CustomerFinalizer        string = "jiraservicedesk.stakater.com/customer"
	CustomerAlreadyExistsErr string = "An account already exists for this email"
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
		log.Info("Adding finalizer for instance " + req.Name)

		finalizerUtil.AddFinalizer(instance, CustomerFinalizer)

		err := r.Client.Update(context.TODO(), instance)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
	}

	// If CustomerId exists in status, then it's an update request
	if len(instance.Status.CustomerId) > 0 {
		// Get the customer from Jira Service Desk
		existingCustomer, err := r.JiraServiceDeskClient.GetCustomerById(instance.Status.CustomerId)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}

		// Check if the customer needs an update
		if r.JiraServiceDeskClient.IsCustomerUpdated(instance, existingCustomer) {

			// Check if this is a valid customer update
			existingCustomerInstance := r.JiraServiceDeskClient.GetCustomerCRFromCustomer(existingCustomer)
			if ok, err := instance.IsValidCustomerUpdate(existingCustomerInstance); !ok {
				return reconcilerUtil.ManageError(r.Client, instance, err, false)
			}

			// Handle customer update
			return r.handleUpdate(req, instance)
		} else {
			log.Info("Skipping update. No changes found")
			return reconcilerUtil.DoNotRequeue()
		}
	}

	return r.handleCreate(req, instance)
}

func (r *CustomerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jiraservicedeskv1alpha1.Customer{}).
		Complete(r)
}

func (r *CustomerReconciler) handleUpdate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	log.Info("Modifying project associations for JSD Customer: " + instance.Spec.Name)

	for _, specProjectKey := range instance.Spec.Projects {
		found := false
		for _, statusProjectKey := range instance.Status.AssociatedProjects {
			if specProjectKey == statusProjectKey {
				found = true
				break
			}
		}
		if !found {
			err := r.JiraServiceDeskClient.AddCustomerToProject(instance.Status.CustomerId, specProjectKey)
			if err != nil {
				return reconcilerUtil.ManageError(r.Client, instance, err, false)
			}
			log.Info("Successfully added Jira Service Desk Customer into project: " + specProjectKey)
		}
	}

	for _, statusProjectKey := range instance.Status.AssociatedProjects {
		found := false
		for _, specProjectKey := range instance.Spec.Projects {
			if specProjectKey == statusProjectKey {
				found = true
				break
			}
		}
		if !found {
			err := r.JiraServiceDeskClient.RemoveCustomerFromProject(instance.Status.CustomerId, statusProjectKey)
			if err != nil {
				return reconcilerUtil.ManageError(r.Client, instance, err, false)
			}
			log.Info("Successfully removed Jira Service Desk Customer from project: " + statusProjectKey)
		}
	}

	instance.Status.AssociatedProjects = instance.Spec.Projects

	return reconcilerUtil.ManageSuccess(r.Client, instance)
}

func (r *CustomerReconciler) handleCreate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	log.Info("Creating Jira Service Desk Customer: " + instance.Spec.Name)
	var customerID string
	var err error

	// If legacy Customer flag is true than create a legacy customer, else create a normal customer
	if instance.Spec.LegacyCustomer {
		customerID, err = r.JiraServiceDeskClient.CreateLegacyCustomer(instance.Spec.Email, instance.Spec.Projects[0])
	} else {
		customer := r.JiraServiceDeskClient.GetCustomerFromCustomerCRForCreateCustomer(instance)
		customerID, err = r.JiraServiceDeskClient.CreateCustomer(customer)
	}

	// If customer already exists, reconstruct status of the custom resource
	if err != nil && strings.Contains(err.Error(), CustomerAlreadyExistsErr) {
		existingCustomerID, err := r.JiraServiceDeskClient.GetCustomerIdByEmail(instance.Spec.Email)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		customerID = existingCustomerID

	} else if err != nil {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	instance.Status.CustomerId = customerID

	log.Info("Successfully created Jira Service Desk Customer: " + instance.Spec.Name)

	log.Info("Adding project associations for JSD Customer: " + instance.Spec.Name)

	for _, projectKey := range instance.Spec.Projects {
		err := r.JiraServiceDeskClient.AddCustomerToProject(instance.Status.CustomerId, projectKey)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		log.Info("Successfully added Jira Service Desk Customer into project: " + projectKey)
	}
	instance.Status.AssociatedProjects = instance.Spec.Projects

	return reconcilerUtil.ManageSuccess(r.Client, instance)
}

func (r *CustomerReconciler) handleDelete(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	if instance == nil {
		// Instance not found, nothing to do
		return reconcilerUtil.DoNotRequeue()
	}

	log.Info("Deleting Jira Service Desk Customer: " + instance.Spec.Name)

	// Check if the customer was created
	if instance.Status.CustomerId != "" {
		// Delete Customer
		err := r.JiraServiceDeskClient.DeleteCustomer(instance.Status.CustomerId)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
	} else {
		log.Info("Customer '" + instance.Spec.Name + "' do not exists on JSD. So skipping deletion")
	}

	// Delete Finalizer
	finalizerUtil.DeleteFinalizer(instance, CustomerFinalizer)

	log.Info("Finalizer removed for customer: " + instance.Spec.Name)

	// Update instance
	err := r.Client.Update(context.TODO(), instance)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	return reconcilerUtil.DoNotRequeue()
}
