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

// RunbookSpec defines the desired state of Runbook
type RunbookSpec struct {
	// AlertName is the name of the associated Prometheus alert
	// +kubebuilder:validation:Required
	AlertName string `json:"alertName"`

	// Severity indicates the alert severity level
	// +kubebuilder:validation:Enum=critical;warning;info
	// +kubebuilder:default=warning
	Severity string `json:"severity,omitempty"`

	// Team responsible for this runbook
	Team string `json:"team,omitempty"`

	// Content contains the runbook documentation
	Content RunbookContent `json:"content"`

	// Template specifies which template to use for generation
	Template string `json:"template,omitempty"`

	// AutoGenerate indicates if this runbook should be auto-generated
	// +kubebuilder:default=true
	AutoGenerate bool `json:"autoGenerate,omitempty"`

	// Outputs specifies where the runbook should be published
	Outputs []OutputConfig `json:"outputs,omitempty"`
}

// RunbookContent defines the structure of runbook documentation
type RunbookContent struct {
	// Impact describes what systems/users are affected
	Impact string `json:"impact,omitempty"`

	// Investigation steps to diagnose the issue
	Investigation []InvestigationStep `json:"investigation,omitempty"`

	// Remediation steps to resolve the issue
	Remediation []RemediationStep `json:"remediation,omitempty"`

	// Prevention describes how to prevent this issue
	Prevention string `json:"prevention,omitempty"`

	// Automation configuration for automatic remediation
	Automation *AutomationConfig `json:"automation,omitempty"`

	// References to external documentation
	References []Reference `json:"references,omitempty"`
}

// InvestigationStep represents a single investigation step
type InvestigationStep struct {
	// Description of what to investigate
	Description string `json:"description"`

	// Command to execute (optional)
	Command string `json:"command,omitempty"`

	// Expected result or what to look for
	Expected string `json:"expected,omitempty"`
}

// RemediationStep represents a single remediation action
type RemediationStep struct {
	// Description of the remediation action
	Description string `json:"description"`

	// Command to execute (optional)
	Command string `json:"command,omitempty"`

	// Risk level of this action
	// +kubebuilder:validation:Enum=low;medium;high
	Risk string `json:"risk,omitempty"`

	// Whether this step can be automated
	Automated bool `json:"automated,omitempty"`
}

// AutomationConfig defines automation settings
type AutomationConfig struct {
	// Enabled indicates if automation is enabled
	Enabled bool `json:"enabled,omitempty"`

	// Scripts to execute for automatic remediation
	Scripts []string `json:"scripts,omitempty"`

	// Triggers define when automation should run
	Triggers []TriggerConfig `json:"triggers,omitempty"`
}

// TriggerConfig defines when automation should trigger
type TriggerConfig struct {
	// Type of trigger (alert, webhook, manual)
	// +kubebuilder:validation:Enum=alert;webhook;manual
	Type string `json:"type"`

	// Conditions that must be met
	Conditions []string `json:"conditions,omitempty"`
}

// Reference represents external documentation links
type Reference struct {
	// Title of the reference
	Title string `json:"title"`

	// URL to the reference
	URL string `json:"url"`

	// Type of reference (wiki, dashboard, documentation)
	// +kubebuilder:validation:Enum=wiki;dashboard;documentation;runbook
	Type string `json:"type,omitempty"`
}

// OutputConfig defines where runbooks should be published
type OutputConfig struct {
	// Format of the output (markdown, html, pdf)
	// +kubebuilder:validation:Enum=markdown;html;pdf
	Format string `json:"format"`

	// Destination where the output should be published
	Destination string `json:"destination"`

	// Template to use for this output
	Template string `json:"template,omitempty"`
}

// RunbookStatus defines the observed state of Runbook
type RunbookStatus struct {
	// Phase represents the current phase of the runbook
	// +kubebuilder:validation:Enum=pending;generating;ready;error
	Phase string `json:"phase,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastGenerated timestamp of last generation
	LastGenerated *metav1.Time `json:"lastGenerated,omitempty"`

	// ValidationStatus indicates if the runbook is valid
	// +kubebuilder:validation:Enum=valid;invalid;pending
	ValidationStatus string `json:"validationStatus,omitempty"`

	// ValidationErrors contains validation error messages
	ValidationErrors []string `json:"validationErrors,omitempty"`

	// GeneratedOutputs lists the successfully generated outputs
	GeneratedOutputs []GeneratedOutput `json:"generatedOutputs,omitempty"`

	// SourceRule reference to the PrometheusRule that generated this runbook
	SourceRule *SourceRuleRef `json:"sourceRule,omitempty"`
}

// GeneratedOutput represents a successfully generated output
type GeneratedOutput struct {
	// Format of the generated output
	Format string `json:"format"`

	// Location where the output was published
	Location string `json:"location"`

	// Timestamp when this output was generated
	GeneratedAt metav1.Time `json:"generatedAt"`
}

// SourceRuleRef references the source PrometheusRule
type SourceRuleRef struct {
	// Name of the PrometheusRule
	Name string `json:"name"`

	// Namespace of the PrometheusRule
	Namespace string `json:"namespace"`

	// UID of the PrometheusRule
	UID string `json:"uid,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced
//+kubebuilder:printcolumn:name="Alert",type=string,JSONPath=`.spec.alertName`
//+kubebuilder:printcolumn:name="Severity",type=string,JSONPath=`.spec.severity`
//+kubebuilder:printcolumn:name="Team",type=string,JSONPath=`.spec.team`
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Runbook is the Schema for the runbooks API
type Runbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RunbookSpec   `json:"spec,omitempty"`
	Status RunbookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RunbookList contains a list of Runbook
type RunbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Runbook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Runbook{}, &RunbookList{})
}
