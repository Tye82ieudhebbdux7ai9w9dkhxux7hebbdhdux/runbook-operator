/*
Copyright 2025 Geovane Guibes.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RunbookTemplateSpec defines the desired state of RunbookTemplate
type RunbookTemplateSpec struct {
	// Name of the template
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description of what this template is for
	Description string `json:"description,omitempty"`

	// Template content using Go template syntax
	// +kubebuilder:validation:Required
	Template string `json:"template"`

	// Variables that can be used in the template
	Variables map[string]TemplateVariable `json:"variables,omitempty"`

	// OutputFormats supported by this template
	OutputFormats []string `json:"outputFormats,omitempty"`

	// Metadata for template organization
	Metadata TemplateMetadata `json:"metadata,omitempty"`
}

// TemplateVariable defines a variable that can be used in templates
type TemplateVariable struct {
	// Description of the variable
	Description string `json:"description,omitempty"`

	// Type of the variable (string, number, boolean, array)
	// +kubebuilder:validation:Enum=string;number;boolean;array
	Type string `json:"type,omitempty"`

	// Default value for the variable
	Default string `json:"default,omitempty"`

	// Whether this variable is required
	Required bool `json:"required,omitempty"`
}

// TemplateMetadata contains metadata for template organization
type TemplateMetadata struct {
	// Tags for categorizing templates
	Tags []string `json:"tags,omitempty"`

	// Author of the template
	Author string `json:"author,omitempty"`

	// Version of the template
	Version string `json:"version,omitempty"`

	// Team responsible for maintaining this template
	Team string `json:"team,omitempty"`
}

// RunbookTemplateStatus defines the observed state of RunbookTemplate
type RunbookTemplateStatus struct {
	// Phase represents the current phase of the template
	// +kubebuilder:validation:Enum=pending;ready;error
	Phase string `json:"phase,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ValidationStatus indicates if the template is valid
	// +kubebuilder:validation:Enum=valid;invalid;pending
	ValidationStatus string `json:"validationStatus,omitempty"`

	// ValidationErrors contains validation error messages
	ValidationErrors []string `json:"validationErrors,omitempty"`

	// UsageCount tracks how many runbooks use this template
	UsageCount int `json:"usageCount,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Template",type=string,JSONPath=`.spec.name`
//+kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.metadata.version`
//+kubebuilder:printcolumn:name="Author",type=string,JSONPath=`.spec.metadata.author`
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// RunbookTemplate is the Schema for the runbooktemplates API
type RunbookTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RunbookTemplateSpec   `json:"spec,omitempty"`
	Status RunbookTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RunbookTemplateList contains a list of RunbookTemplate
type RunbookTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RunbookTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RunbookTemplate{}, &RunbookTemplateList{})
}
