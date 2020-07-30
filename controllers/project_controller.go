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
	// 	project, err := r.JiraServiceDeskClient.GetProjectByKey(instance.Spec.Key)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Project already exists
	// 	// TODO: This should be project != nil
	// 	if err != nil {
	// 		updatedProject := r.JiraServiceDeskClient.GetProjectFromProjectSpec(instance.Spec)
	// 		if !r.JiraServiceDeskClient.ProjectEqual(project, updatedProject) {
	// 			return r.handleUpdate(req, instance)
	// 		} else {
	// 			log.Info("Skipping update. No changes found")
	// 			return ctrl.Result{}, nil
	// 		}
	// 	}
	// TODO: Think of use cases and add a default return ctrl.Result{}, nil
	return r.handleCreate(req, instance, log)
}

func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jiraservicedeskv1alpha1.Project{}).
		Complete(r)
}

func (r *ProjectReconciler) handleCreate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project, log logr.Logger) (ctrl.Result, error) {

	log.Info("Creating Jira Service Desk Project: " + instance.Spec.Name)

	project := r.JiraServiceDeskClient.GetProjectFromProjectSpec(instance.Spec)
	err := r.JiraServiceDeskClient.CreateProject(project)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Successfully created Jira Service Desk Project: " + instance.Spec.Name)

	return ctrl.Result{RequeueAfter: defaultRequeueTime}, nil
}

func (r *ProjectReconciler) handleDelete(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *ProjectReconciler) handleUpdate(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: defaultRequeueTime}, nil
}
