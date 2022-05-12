/*
Copyright 2021.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var customerlog = logf.Log.WithName("customer-resource")

func (r *Customer) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-jiraservicedesk-stakater-com-v1alpha1-customer,mutating=true,failurePolicy=fail,sideEffects=None,groups=jiraservicedesk.stakater.com,resources=customers,verbs=create;update,versions=v1alpha1,name=mcustomer.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Customer{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Customer) Default() {
	customerlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-jiraservicedesk-stakater-com-v1alpha1-customer,mutating=false,failurePolicy=fail,sideEffects=None,groups=jiraservicedesk.stakater.com,resources=customers,verbs=create;update,versions=v1alpha1,name=vcustomer.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Customer{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Customer) ValidateCreate() error {
	customerlog.Info("validate create", "name", r.Name)

	_, err := r.IsValid()
	return err
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Customer) ValidateUpdate(old runtime.Object) error {
	customerlog.Info("validate update", "name", r.Name)

	oldCustomer, ok := old.(*Customer)
	if !ok {
		return fmt.Errorf("Error casting old runtime object to %T from %T", oldCustomer, old)
	}

	_, err := r.IsValid()
	if err != nil {
		return err
	}
	_, err = r.IsValidUpdate(*oldCustomer)
	if err != nil {
		return err
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Customer) ValidateDelete() error {
	customerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
