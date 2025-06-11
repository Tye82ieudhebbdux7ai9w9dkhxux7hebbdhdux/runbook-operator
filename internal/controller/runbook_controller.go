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

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&runbook, "runbook.runbook.io/finalizer") {
		controllerutil.AddFinalizer(&runbook, "runbook.runbook.io/finalizer")
		if err := r.Update(ctx, &runbook); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Handle deletion
	if runbook.DeletionTimestamp != nil {
		return r.handleDeletion(ctx, &runbook)
	}

	// Reconcile the runbook
	return r.reconcileRunbook(ctx, &runbook)
}

func (r *RunbookReconciler) reconcileRunbook(ctx context.Context, runbook *runbookv1alpha1.Runbook) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Update status to generating
	if runbook.Status.Phase != "generating" {
		runbook.Status.Phase = "generating"
		if err := r.Status().Update(ctx, runbook); err != nil {
			return ctrl.Result{}, err
		}
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

	if err := r.Status().Update(ctx, runbook); err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled runbook", "runbook", runbook.Name)
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

func (r *RunbookReconciler) generateRunbookContent(ctx context.Context, runbook *runbookv1alpha1.Runbook) error {
	// TODO: Implement runbook content generation
	// This will integrate with the generator package
	return nil
}

func (r *RunbookReconciler) validateRunbook(ctx context.Context, runbook *runbookv1alpha1.Runbook) error {
	// TODO: Implement runbook validation
	// This will integrate with the validator package
	return nil
}

func (r *RunbookReconciler) generateOutputs(ctx context.Context, runbook *runbookv1alpha1.Runbook) error {
	// TODO: Implement output generation
	// This will generate markdown, HTML, PDF files
	return nil
}

func (r *RunbookReconciler) updateStatusWithError(ctx context.Context, runbook *runbookv1alpha1.Runbook, err error) (ctrl.Result, error) {
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

	if updateErr := r.Status().Update(ctx, runbook); updateErr != nil {
		return ctrl.Result{}, updateErr
	}

	return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
}

func (r *RunbookReconciler) handleDeletion(ctx context.Context, runbook *runbookv1alpha1.Runbook) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO: Cleanup generated files and external resources

	// Remove finalizer
	controllerutil.RemoveFinalizer(runbook, "runbook.runbook.io/finalizer")
	if err := r.Update(ctx, runbook); err != nil {
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
