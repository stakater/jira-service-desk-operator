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

	if len(instance.Status.CustomerId) > 0 {
		return r.handleUpdate(req, instance)
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

	log.Info("Modifying project associations for jsd customer: " + instance.Spec.Name)

	addedProjects := difference(instance.Spec.Projects, instance.Status.AssociatedProjects)
	removedProjects := difference(instance.Status.AssociatedProjects, instance.Spec.Projects)

	for _, projectKey := range addedProjects {
		err := r.JiraServiceDeskClient.AddCustomerToProject(instance.Status.CustomerId, projectKey)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		instance.Status.AssociatedProjects = append(instance.Status.AssociatedProjects, projectKey)
		log.Info("Successfully added Jira Service Desk Customer into project: " + projectKey)
	}

	for index, projectKey := range removedProjects {
		err := r.JiraServiceDeskClient.RemoveCustomerFromProject(instance.Status.CustomerId, projectKey)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		instance.Status.AssociatedProjects[index] = ""
		log.Info("Successfully removed Jira Service Desk Customer from project: " + projectKey)
	}

	instance.Status.AssociatedProjects = removeEmptyProjects(instance.Status.AssociatedProjects)

	return reconcilerUtil.ManageSuccess(r.Client, instance)
}

func (r *CustomerReconciler) handleCreate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	log.Info("Creating Jira Service Desk Customer: " + instance.Spec.Name)

	customer := r.JiraServiceDeskClient.GetCustomerFromCustomerCRForCreateCustomer(instance)
	customerId, err := r.JiraServiceDeskClient.CreateCustomer(customer)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	instance.Status.CustomerId = customerId
	log.Info("Successfully created Jira Service Desk Customer: " + instance.Spec.Name)

	log.Info("Modifying project associations for jsd customer: " + instance.Spec.Name)

	for _, projectKey := range instance.Spec.Projects {
		err := r.JiraServiceDeskClient.AddCustomerToProject(instance.Status.CustomerId, projectKey)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		instance.Status.AssociatedProjects = append(instance.Status.AssociatedProjects, projectKey)
		log.Info("Successfully added Jira Service Desk Customer into project: " + projectKey)
	}

	instance.Status.AssociatedProjects = removeEmptyProjects(instance.Status.AssociatedProjects)

	return reconcilerUtil.ManageSuccess(r.Client, instance)
}

func (r *CustomerReconciler) handleDelete(req ctrl.Request, instance *jiraservicedeskv1alpha1.Customer) (ctrl.Result, error) {
	log := r.Log.WithValues("customer", req.NamespacedName)

	if instance == nil {
		// Instance not found, nothing to do
		return reconcilerUtil.DoNotRequeue()
	}

	log.Info("Removing project associations for jsd customer: " + instance.Spec.Name)

	// Remove customer from JSD project
	for _, projectKey := range instance.Status.AssociatedProjects {
		err := r.JiraServiceDeskClient.RemoveCustomerFromProject(instance.Status.CustomerId, projectKey)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		log.Info("Successfully removed Jira Service Desk Customer from project: " + projectKey)
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

func removeEmptyProjects(slice []string) []string {
	var output []string
	for _, str := range slice {
		if str != "" {
			output = append(output, str)
		}
	}
	return output
}

// returns the difference between the two slices
func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, obj1 := range slice1 {
		found := false
		for _, obj2 := range slice2 {
			if obj1 == obj2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, obj1)
		}
	}
	return diff
}
