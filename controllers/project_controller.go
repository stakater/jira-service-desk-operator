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
	// TODO: Check if this is required in our case
	// 	defaultRequeueTime        = 60 * time.Second
	ProjectFinalizer        string = "jiraservicedesk.stakater.com/project"
	ProjectAlreadyExistsErr string = "A project with that name already exists."
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
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list

func (r *ProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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
			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading the object - requeue the request.
		return reconcilerUtil.RequeueWithError(err)
	}

	// Validate Custom Resource
	if ok, err := instance.IsValid(); !ok {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	// Resource is marked for deletion
	if instance.DeletionTimestamp != nil {
		log.Info("Deletion timestamp found for instance " + req.Name)
		if finalizerUtil.HasFinalizer(instance, ProjectFinalizer) {
			return r.handleDelete(req, instance)
		}
		// Finalizer doesn't exist so clean up is already done
		return reconcilerUtil.DoNotRequeue()
	}

	// Add finalizer if it doesn't exist
	if !finalizerUtil.HasFinalizer(instance, ProjectFinalizer) {
		log.Info("Adding finalizer for instance " + req.Name)

		finalizerUtil.AddFinalizer(instance, ProjectFinalizer)

		err := r.Client.Update(context.TODO(), instance)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
	}

	// Check if the Project already exists
	if len(instance.Status.ID) > 0 {
		existingProject, err := r.JiraServiceDeskClient.GetProjectByIdentifier(instance.Status.ID)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		// Project already exists
		if len(existingProject.Id) > 0 {
			updatedProject := r.JiraServiceDeskClient.GetProjectFromProjectCR(instance)
			// Compare retrieved project with current spec
			if !r.JiraServiceDeskClient.ProjectEqual(existingProject, updatedProject) {
				// Update if there are changes in the declared spec
				return r.handleUpdate(req, existingProject, instance)
			} else {
				log.Info("Skipping update. No changes found")
				return reconcilerUtil.DoNotRequeue()
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

	project := r.JiraServiceDeskClient.GetProjectFromProjectCR(instance)
	projectId, err := r.JiraServiceDeskClient.CreateProject(project)

	// If project already exists then reconstruct status
	if err != nil && strings.Contains(err.Error(), ProjectAlreadyExistsErr) {
		existingProject, err := r.JiraServiceDeskClient.GetProjectByIdentifier(instance.Spec.Key)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
		log.Info("Successfully reconstructed status for Jira Service Desk Project " + instance.Spec.Name)

		projectId = existingProject.Id
	} else if err != nil {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	log.Info("Successfully created Jira Service Desk Project: " + instance.Spec.Name)

	if !instance.Spec.OpenAccess {
		err = r.JiraServiceDeskClient.UpdateProjectAccessPermissions(instance.Spec.OpenAccess, project.Key)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}

		log.Info("Successfully updated the Access Permissions to customer")
	}

	instance.Status.ID = projectId
	return reconcilerUtil.ManageSuccess(r.Client, instance)
}

func (r *ProjectReconciler) handleDelete(req ctrl.Request, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	log := r.Log.WithValues("project", req.NamespacedName)

	if instance == nil {
		// Instance not found, nothing to do
		return reconcilerUtil.DoNotRequeue()
	}

	log.Info("Deleting Jira Service Desk Project: " + instance.Spec.Name)

	// Check if the project was created
	if instance.Status.ID != "" {
		err := r.JiraServiceDeskClient.DeleteProject(instance.Status.ID)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, instance, err, false)
		}
	} else {
		log.Info("Project '" + instance.Spec.Name + "' do not exists on JSD. So skipping deletion")
	}

	// Delete finalizer
	finalizerUtil.DeleteFinalizer(instance, ProjectFinalizer)

	log.Info("Finalizer removed for project: " + instance.Spec.Name)

	// Update instance
	err := r.Client.Update(context.TODO(), instance)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	return reconcilerUtil.DoNotRequeue()
}

func (r *ProjectReconciler) handleUpdate(req ctrl.Request, existingProject jiraservicedeskclient.Project, instance *jiraservicedeskv1alpha1.Project) (ctrl.Result, error) {
	log := r.Log.WithValues("project", req.NamespacedName)

	log.Info("Updating Jira Service Desk Project: " + instance.Spec.Name)

	existingProjectInstance := r.JiraServiceDeskClient.GetProjectCRFromProject(existingProject)
	if ok, err := instance.IsValidUpdate(existingProjectInstance); !ok {
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	updatedProject := r.JiraServiceDeskClient.GetProjectForUpdateRequest(existingProject, instance)
	err := r.JiraServiceDeskClient.UpdateProject(updatedProject, existingProject.Id)
	if err != nil {
		log.Error(err, "Failed to update status of Project")
		return reconcilerUtil.ManageError(r.Client, instance, err, false)
	}

	return reconcilerUtil.ManageSuccess(r.Client, instance)
}
