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
	"github.com/operator-framework/operator-sdk/pkg/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProjectSpec defines the desired state of Project
type ProjectSpec struct {

	// Name of the project
	// +required
	Name string `json:"name,omitempty"`

	// The project key is used as the prefix of your project's issue keys
	// +required
	Key string `json:"key,omitempty"`

	// The project type, which dictates the application-specific feature set
	// +kubebuilder:validation:Enum=business;service_desk;software
	// +required
	ProjectTypeKey string `json:"projectTypeKey,omitempty"`

	// A prebuilt configuration for a project
	// +required
	ProjectTemplateKey string `json:"projectTemplateKey,omitempty"`

	// Description for project
	// +required
	Description string `json:"description,omitempty"`

	// Task assignee type
	// +kubebuilder:validation:Enum=PROJECT_LEAD;UNASSIGNED
	// +required
	AssigneeType string `json:"assigneeType,omitempty"`

	// ID of project lead
	// +kubebuilder:validation:MaxLength=128
	// +required
	LeadAccountId string `json:"leadAccountId,omitempty"`

	// A link to information about this project, such as project documentation
	// +optional
	URL string `json:"url,omitempty"`

	// An integer value for the project's avatar.
	// +optional
	AvatarId int `json:"avatarId,omitempty"`

	// The ID of the issue security scheme for the project, which enables you to control who can and cannot view issues
	// +optional
	IssueSecurityScheme int `json:"issueSecurityScheme,omitempty"`

	// The ID of the permission scheme for the project
	// +optional
	PermissionScheme int `json:"permissionScheme,omitempty"`

	// The ID of the notification scheme for the project
	// +optional
	NotificationScheme int `json:"notificationScheme,omitempty"`

	// The ID of the project's category
	// +optional
	CategoryId int `json:"categoryId,omitempty"`
}

// ProjectStatus defines the observed state of Project
type ProjectStatus struct {
	// Jira service desk project ID
	ID string `json:"id"`

	// Status conditions
	Conditions status.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Project is the Schema for the projects API
// +kubebuilder:resource:path=projects,scope=Cluster
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProjectList contains a list of Project
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Project{}, &ProjectList{})
}

func (project *Project) GetReconcileStatus() status.Conditions {
	return project.Status.Conditions
}

func (project *Project) SetReconcileStatus(reconcileStatus status.Conditions) {
	project.Status.Conditions = reconcileStatus
}

func (project *Project) IsValid() (bool, error) {
	// Add logic for additional validation here
	return true, nil
}

func (project *Project) ValidateForUpdate() (bool, error) {
	// Add logic for additional validation here
	return true, nil
}
