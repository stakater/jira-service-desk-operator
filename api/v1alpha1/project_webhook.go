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

package v1alpha1

import (
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var projectlog = logf.Log.WithName("project-resource")

func (r *Project) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-jiraservicedesk-stakater-com-v1alpha1-project,mutating=true,failurePolicy=fail,groups=jiraservicedesk.stakater.com,resources=projects,verbs=create;update,versions=v1alpha1,name=mproject.kb.io

var _ webhook.Defaulter = &Project{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Project) Default() {
	projectlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-jiraservicedesk-stakater-com-v1alpha1-project,mutating=false,failurePolicy=fail,groups=jiraservicedesk.stakater.com,resources=projects,versions=v1alpha1,name=vproject.kb.io

var _ webhook.Validator = &Project{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Project) ValidateCreate() error {
	projectlog.Info("validate create", "name", r.Name)

	_, err := r.IsValid()
	return err
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Project) ValidateUpdate(old runtime.Object) error {
	projectlog.Info("validate update", "name", r.Name)

	oldProject, ok := old.(*Project)
	if !ok {
		log.Error(ok, "Error in finding the requested project")
	}
	_, err := r.IsValidUpdate(*oldProject)
	return err
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Project) ValidateDelete() error {
	projectlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
