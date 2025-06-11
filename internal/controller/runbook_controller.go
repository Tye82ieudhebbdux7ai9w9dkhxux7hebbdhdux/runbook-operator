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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	runbookv1alpha1 "github.com/guibes/runbook-operator/api/v1alpha1"
	"github.com/guibes/runbook-operator/pkg/generator"
	"github.com/guibes/runbook-operator/pkg/outputs"
)

// RunbookReconciler reconciles a Runbook object
type RunbookReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Generator *generator.RunbookGenerator
}

//+kubebuilder:rbac:groups=runbook.runbook.io,resources=runbooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=runbook.runbook.io,resources=runbooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=runbook.runbook.io,resources=runbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop
func (r *RunbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Runbook instance
	var runbook runbookv1alpha1.Runbook
	if err := r.Get(ctx, req.NamespacedName, &runbook); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Runbook resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Runbook")
		return ctrl.Result{}, err
	}

	logger.Info("🔥 Reconciling runbook", "runbook", runbook.Name)

	// Handle deletion
	if runbook.DeletionTimestamp != nil {
		return r.handleDeletion(ctx, &runbook)
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&runbook, "runbook.runbook.io/finalizer") {
		original := runbook.DeepCopy()
		controllerutil.AddFinalizer(&runbook, "runbook.runbook.io/finalizer")
		if err := r.Patch(ctx, &runbook, client.MergeFrom(original)); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		logger.Info("Added finalizer to runbook")
		return ctrl.Result{Requeue: true}, nil
	}

	// Reconcile the runbook
	return r.reconcileRunbook(ctx, &runbook)
}

func (r *RunbookReconciler) reconcileRunbook(ctx context.Context, runbook *runbookv1alpha1.Runbook) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Create a copy for status updates
	original := runbook.DeepCopy()

	// Update status to generating if not already set
	if runbook.Status.Phase != "generating" && runbook.Status.Phase != "ready" {
		runbook.Status.Phase = "generating"
		if err := r.updateStatus(ctx, runbook, original); err != nil {
			return ctrl.Result{}, err
		}
		// Return and requeue to get the updated object
		return ctrl.Result{Requeue: true}, nil
	}

	// Generate runbook content if auto-generate is enabled
	if runbook.Spec.AutoGenerate {
		if err := r.generateRunbookContent(ctx, runbook); err != nil {
			logger.Error(err, "Failed to generate runbook content")
			return r.updateStatusWithError(ctx, runbook, err)
		}
	}

	// Validate runbook content
	if err := r.validateRunbook(ctx, runbook); err != nil {
		logger.Error(err, "Failed to validate runbook")
		return r.updateStatusWithError(ctx, runbook, err)
	}

	// Generate outputs
	if err := r.generateOutputs(ctx, runbook); err != nil {
		logger.Error(err, "Failed to generate outputs")
		return r.updateStatusWithError(ctx, runbook, err)
	}

	// Update status to ready
	runbook.Status.Phase = "ready"
	runbook.Status.ValidationStatus = "valid"
	runbook.Status.ValidationErrors = nil // Clear any previous errors
	now := metav1.NewTime(time.Now())
	runbook.Status.LastGenerated = &now

	// Update condition
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		Reason:             "GenerationSuccessful",
		Message:            "Runbook generation completed successfully",
	}
	runbook.Status.Conditions = []metav1.Condition{condition}

	if err := r.updateStatus(ctx, runbook, original); err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled runbook", "runbook", runbook.Name)
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// updateStatus safely updates the status with retry logic
func (r *RunbookReconciler) updateStatus(ctx context.Context, runbook *runbookv1alpha1.Runbook, original *runbookv1alpha1.Runbook) error {
	logger := log.FromContext(ctx)

	// Use patch instead of update to avoid conflicts
	if err := r.Status().Patch(ctx, runbook, client.MergeFrom(original)); err != nil {
		if errors.IsConflict(err) {
			logger.Info("Status update conflict, will retry on next reconciliation")
			return nil // Don't return error, let the reconciler retry
		}
		logger.Error(err, "Failed to update status")
		return err
	}
	return nil
}

func (r *RunbookReconciler) generateRunbookContent(ctx context.Context, runbook *runbookv1alpha1.Runbook) error {
	logger := log.FromContext(ctx)
	logger.Info("Generating runbook content", "runbook", runbook.Name)

	// TODO: Implement actual runbook content generation
	// For now, just ensure we have some basic content
	if runbook.Spec.Content.Impact == "" {
		logger.Info("No impact defined, using default")
	}

	return nil
}

func (r *RunbookReconciler) validateRunbook(ctx context.Context, runbook *runbookv1alpha1.Runbook) error {
	logger := log.FromContext(ctx)
	logger.Info("Validating runbook", "runbook", runbook.Name)

	// TODO: Implement actual validation
	// Basic validation for now
	if runbook.Spec.AlertName == "" {
		return fmt.Errorf("alert name is required")
	}

	return nil
}

func (r *RunbookReconciler) generateOutputs(ctx context.Context, runbook *runbookv1alpha1.Runbook) error {
	logger := log.FromContext(ctx)

	content, err := r.Generator.GenerateMarkdown(ctx, runbook)
	if err != nil {
		return fmt.Errorf("failed to generate markdown content: %w", err)
	}

	var generatedOutputs []runbookv1alpha1.GeneratedOutput

	for _, output := range runbook.Spec.Outputs {
		logger.Info("Generating output", "type", output.Format, "runbook", runbook.Name)

		var err error
		switch output.Format {
		case "markdown":
			mardownOut := &outputs.MarkdownOutput{BasePath: output.Destination}
			err = mardownOut.Generate(runbook, content)
		case "html":
			htmlOut := &outputs.HTMLOutput{BasePath: output.Destination}
			err = htmlOut.Generate(runbook)
		case "api":
			apiOut := &outputs.APIOutput{BaseURL: output.Destination}
			err = apiOut.Generate(runbook, content)
		default:
			logger.Info("Unknown output format, skipping", "format", output.Format)
			continue
		}

		if err != nil {
			logger.Error(err, "Failed to generate output", "type", output.Format)
			continue
		}

		generatedOutputs = append(generatedOutputs, runbookv1alpha1.GeneratedOutput{
			Format:      output.Format,
			Location:    output.Destination,
			GeneratedAt: metav1.NewTime(time.Now()),
		})

		runbook.Status.GeneratedOutputs = generatedOutputs
	}

	return nil
}

func (r *RunbookReconciler) updateStatusWithError(ctx context.Context, runbook *runbookv1alpha1.Runbook, err error) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	original := runbook.DeepCopy()

	runbook.Status.Phase = "error"
	runbook.Status.ValidationStatus = "invalid"
	runbook.Status.ValidationErrors = []string{err.Error()}

	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Reason:             "GenerationFailed",
		Message:            fmt.Sprintf("Failed to generate runbook: %v", err),
	}
	runbook.Status.Conditions = []metav1.Condition{condition}

	if updateErr := r.updateStatus(ctx, runbook, original); updateErr != nil {
		logger.Error(updateErr, "Failed to update error status")
	}

	return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
}

func (r *RunbookReconciler) handleDeletion(ctx context.Context, runbook *runbookv1alpha1.Runbook) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO: Cleanup generated files and external resources
	logger.Info("Cleaning up runbook resources", "runbook", runbook.Name)

	// Remove finalizer
	original := runbook.DeepCopy()
	controllerutil.RemoveFinalizer(runbook, "runbook.runbook.io/finalizer")
	if err := r.Patch(ctx, runbook, client.MergeFrom(original)); err != nil {
		if errors.IsConflict(err) {
			logger.Info("Conflict while removing finalizer, will retry")
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	logger.Info("Successfully deleted runbook", "runbook", runbook.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RunbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runbookv1alpha1.Runbook{}).
		Complete(r)
}
