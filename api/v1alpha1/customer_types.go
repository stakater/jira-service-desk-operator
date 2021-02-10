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
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	invalidUpdateErrorMsg string = " is an immutable field and can not be modified."
	duplicateKeysErr      string = "Duplicate Project Keys are not allowed"
)

// CustomerSpec defines the desired state of Customer
type CustomerSpec struct {
	// Name of the customer
	// +required
	Name string `json:"name,omitempty"`

	// Email of the customer
	// +kubebuilder:validation:Pattern=\S+@\S+\.\S+
	// +required
	Email string `json:"email,omitempty"`

	// LegacyCustomer is a boolean flag that represents whether a customer is created using legacy API or not
	// In case of a legacy Customer, a signup link is sent to the customer email which he can than use to signup
	// In case of a normal Customer, no signup link is sent to the customer. The customer than has to signup manually using the portal
	// If not given, default behaviour is false i.e. normal customer
	// +optional
	LegacyCustomer bool `json:"legacyCustomer,omitempty"`

	// List of ProjectKeys in which customer will be added
	// +kubebuilder:validation:MinItems=1
	// +optional
	Projects []string `json:"projects,omitempty"`
}

// CustomerStatus defines the observed state of Customer
type CustomerStatus struct {
	// Jira Service Desk Customer Account Id
	CustomerId string `json:"customerId"`

	// List of ProjectKeys in which customer has bee added
	AssociatedProjects []string `json:"associatedProjects,omitempty"`

	// Status conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Customer is the Schema for the customers API
type Customer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomerSpec   `json:"spec,omitempty"`
	Status CustomerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomerList contains a list of Customer
type CustomerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Customer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Customer{}, &CustomerList{})
}

func (customer *Customer) GetReconcileStatus() []metav1.Condition {
	return customer.Status.Conditions
}

func (customer *Customer) SetReconcileStatus(reconcileStatus []metav1.Condition) {
	customer.Status.Conditions = reconcileStatus
}

func (customer *Customer) IsValid() (bool, error) {

	if duplicateKeysExist(customer.Spec.Projects) {
		return false, errors.New(duplicateKeysErr)
	}

	return true, nil
}

func (customer *Customer) IsValidUpdate(existingCustomer Customer) (bool, error) {

	if !strings.EqualFold(customer.Spec.Email, existingCustomer.Spec.Email) {
		return false, fmt.Errorf("%s %s", "CustomerEmail", invalidUpdateErrorMsg)
	}
	if customer.Spec.Name != existingCustomer.Spec.Name {
		return false, fmt.Errorf("%s %s", "CustomerName", invalidUpdateErrorMsg)
	}

	if customer.Spec.LegacyCustomer != existingCustomer.Spec.LegacyCustomer {
		return false, fmt.Errorf("%s %s", "LegacyCustomer", invalidUpdateErrorMsg)
	}

	if duplicateKeysExist(customer.Spec.Projects) {
		return false, errors.New(duplicateKeysErr)
	}

	return true, nil
}

func (customer *Customer) IsValidCustomerUpdate(existingCustomer Customer) (bool, error) {
	if !strings.EqualFold(customer.Spec.Email, existingCustomer.Spec.Email) {
		// It takes a few seconds for customers to be persisted at JSD. Check if it's pending with known values are return.
		if existingCustomer.Spec.Name == "User " && existingCustomer.Spec.Email == "?" {
			return true, nil
		}
		return false, fmt.Errorf("%s %s", "Customer email", invalidUpdateErrorMsg)
	}

	return true, nil
}

func duplicateKeysExist(projectKeys []string) bool {
	keys := make(map[string]bool)

	for _, entry := range projectKeys {
		if _, value := keys[entry]; !value {
			keys[entry] = true
		} else {
			return true
		}
	}
	return false
}
