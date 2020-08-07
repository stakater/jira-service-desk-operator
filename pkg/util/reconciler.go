package util

import (
	"context"
	"time"

	astatus "github.com/operator-framework/operator-sdk/pkg/ansible/controller/status"
	"github.com/operator-framework/operator-sdk/pkg/status"
	apis "github.com/redhat-cop/operator-utils/pkg/util/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("util")

//ManageError will set status of the passed CR to a error condition
func ManageError(client client.Client, obj apis.Resource, issue error) (ctrl.Result, error) {
	if reconcileStatusAware, updateStatus := (obj).(apis.ConditionsStatusAware); updateStatus {
		condition := status.Condition{
			Type:               "ReconcileError",
			LastTransitionTime: metav1.Now(),
			Message:            issue.Error(),
			Reason:             astatus.FailedReason,
			Status:             corev1.ConditionTrue,
		}
		reconcileStatusAware.SetReconcileStatus(status.NewConditions(condition))
		err := client.Status().Update(context.Background(), obj)
		if err != nil {
			log.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
	} else {
		log.Info("object is not ReconcileStatusAware, not setting status")
	}
	return ctrl.Result{}, issue
}

// ManageSuccess will update the status of the CR and return a successful reconcile result
func ManageSuccess(client client.Client, obj apis.Resource) (ctrl.Result, error) {
	if reconcileStatusAware, updateStatus := (obj).(apis.ConditionsStatusAware); updateStatus {
		condition := status.Condition{
			Type:               "ReconcileSuccess",
			LastTransitionTime: metav1.Now(),
			Message:            astatus.SuccessfulMessage,
			Reason:             astatus.SuccessfulReason,
			Status:             corev1.ConditionTrue,
		}
		reconcileStatusAware.SetReconcileStatus(status.NewConditions(condition))
		err := client.Status().Update(context.Background(), obj)
		if err != nil {
			log.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
	} else {
		log.Info("object is not ReconcileStatusAware, not setting status")
	}
	return ctrl.Result{}, nil
}

func DoNotRequeue() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func RequeueWithError(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func RequeueAfter(requeueTime time.Duration) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: requeueTime}, nil
}
