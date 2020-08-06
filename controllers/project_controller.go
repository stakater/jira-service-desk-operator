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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
	jiraservicedeskclient "github.com/stakater/jira-service-desk-operator/jiraservicedesk/client"
)

const (
	defaultRequeueTime = 60 * time.Second
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Scheme                *runtime.Scheme
	Log                   logr.Logger
	JiraServiceDeskClient jiraservicedeskclient.Client
}

// +kubebuilder:rbac:groups=jiraservicedesk.stakater.com,resources=projects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jiraservicedesk.stakater.com,resources=projects/status,verbs=get;update;patch

func (r *ProjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	log := r.Log.WithValues("project", req.NamespacedName)

	log.Info("Reconciling Project")

	// Fetch the Project instance
	instance := &jiraservicedeskv1alpha1.Project{}

	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return r.handleDelete(req, instance)
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Check if the Project already exists
	if len(instance.Status.ID) > 0 {
		project, err := r.JiraServiceDeskClient.GetProjectById(instance.Status.ID)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Project already exists
		if len(project.Id) > 0 {
			updatedProject := r.JiraServiceDeskClient.GetProjectFromProjectSpec(instance.Spec)
			// Compare retrieved project with current spec
			if !r.JiraServiceDeskClient.ProjectEqual(project, updatedProject) {
				// Update if there are changes in the declared spec
				return r.handleUpdate(req, instance)
			} else {
				log.Info("Skipping update. No changes found")
				return ctrl.Result{}, nil
			}
		}
	}
	return r.handleCreate(req, instance)
}

func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jiraservicedeskv1alpha1.Project{}).
		Complete(r)
}

func (r *ProjectReconciler) handleCreate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	log := r.Log.WithValues("project", req.NamespacedName)

	log.Info("Creating Jira Service Desk Project: " + instance.Spec.Name)

	project := r.JiraServiceDeskClient.GetProjectFromProjectSpec(instance.Spec)
	projectId, err := r.JiraServiceDeskClient.CreateProject(project)
	if err != nil {
		return ctrl.Result{}, err
	}

	instance.Status.ID = projectId

	err = r.Status().Update(context.Background(), instance)
	if err != nil {
		log.Error(err, "Failed to update status of Project")
		return ctrl.Result{}, err
	}

	log.Info("Successfully created Jira Service Desk Project: " + instance.Spec.Name)

	return ctrl.Result{RequeueAfter: defaultRequeueTime}, nil
}

func (r *ProjectReconciler) handleDelete(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	log := r.Log.WithValues("project", req.NamespacedName)

	log.Info("Deleting Jira Service Desk Project: " + instance.Spec.Name)

	if instance == nil {
		// Instance not found, nothing to do
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *ProjectReconciler) handleUpdate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	log := r.Log.WithValues("project", req.NamespacedName)

	log.Info("Updating Jira Service Desk Project: " + instance.Spec.Name)

	return ctrl.Result{RequeueAfter: defaultRequeueTime}, nil
}
