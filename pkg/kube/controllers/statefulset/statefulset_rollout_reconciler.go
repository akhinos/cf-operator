package statefulset

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/cf-operator/pkg/kube/apis"

	"code.cloudfoundry.org/quarks-utils/pkg/config"
	"code.cloudfoundry.org/quarks-utils/pkg/ctxlog"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	crc "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	rolloutStatePending       = "Pending"
	rolloutStateCanary        = "Canary"
	rolloutStateRollout       = "Rollout"
	rolloutStateDone          = "Done"
	rolloutStateFailed        = "Failed"
	rolloutStateCanaryUpscale = "CanaryUpscale"
)

var (
	// annotationCanaryRollout is the state of the canary rollout of the stateful set
	annotationCanaryRollout = fmt.Sprintf("%s/canary-rollout", apis.GroupName)
)

// NewStatefulSetRolloutReconciler returns a new reconcile.Reconciler
func NewStatefulSetRolloutReconciler(ctx context.Context, config *config.Config, mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStatefulSetRollout{
		ctx:    ctx,
		config: config,
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// ReconcileStatefulSetRollout reconciles an ExtendedStatefulSet object when references changes
type ReconcileStatefulSetRollout struct {
	ctx    context.Context
	client crc.Client
	scheme *runtime.Scheme
	config *config.Config
}

// Reconcile cleans up old versions and volumeManagement statefulSet of the ExtendedStatefulSet
func (r *ReconcileStatefulSetRollout) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Set the ctx to be Background, as the top-level context for incoming requests.
	ctx, cancel := context.WithTimeout(r.ctx, r.config.CtxTimeOut)
	defer cancel()

	ctxlog.Debug(ctx, "Reconciling StatefulSet ", request.NamespacedName)

	statefulSet := v1beta2.StatefulSet{}

	err := r.client.Get(ctx, request.NamespacedName, &statefulSet)
	if err != nil {
		ctxlog.Debug(ctx, "StatefulSet not found ", request.NamespacedName)
		return reconcile.Result{}, err
	}

	dirty := false

	var status = r.getState(&statefulSet)
	var newStatus = status
	switch status {
	case rolloutStateCanaryUpscale:
		if statefulSet.Status.Replicas == *statefulSet.Spec.Replicas && statefulSet.Status.ReadyReplicas == *statefulSet.Spec.Replicas {
			if *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition == 0 {
				newStatus = rolloutStateDone
			} else {
				newStatus = rolloutStateRollout
			}
		}
		break
	case rolloutStateDone:
		break
	case rolloutStateFailed:
		break
	case rolloutStateCanary:
		fallthrough
	case rolloutStateRollout:
		ready, err := readyAndUpdated(ctx, r.client, &statefulSet)
		if err != nil {
			return reconcile.Result{}, err
		}
		if *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition == 0 {
			if ready {
				newStatus = rolloutStateDone
			}
			break
		}
		if !ready {
			break
		}
		(*statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition)--
		dirty = true
		err = CleanupNonReadyPod(ctx, r.client, &statefulSet, *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition)
		if err != nil {
			ctxlog.Debug(ctx, "Error calling CleanupNonReadyPod ", request.NamespacedName, err)
			return reconcile.Result{}, err
		}
		newStatus = r.getState(&statefulSet)
		break
	case rolloutStatePending:
		if statefulSet.Status.Replicas < *statefulSet.Spec.Replicas {
			newStatus = rolloutStateCanaryUpscale
		} else {
			newStatus = rolloutStateCanary
			err = CleanupNonReadyPod(ctx, r.client, &statefulSet, *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition)
			if err != nil {
				ctxlog.Debug(ctx, "Error calling CleanupNonReadyPod ", request.NamespacedName, err)
				return reconcile.Result{}, err
			}
		}
		break
	}
	if newStatus != statefulSet.Annotations[annotationCanaryRollout] {
		statefulSet.Annotations[annotationCanaryRollout] = newStatus
		dirty = true
	}
	if dirty {
		r.update(ctx, &statefulSet)
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileStatefulSetRollout) getState(sts *v1beta2.StatefulSet) string {
	state := sts.Annotations[annotationCanaryRollout]
	switch state {
	case rolloutStateFailed:
		return state
	case rolloutStatePending:
		return state
	case rolloutStateCanaryUpscale:
		return state
	default:
		if *sts.Spec.UpdateStrategy.RollingUpdate.Partition == *sts.Spec.Replicas-1 {
			return rolloutStateCanary
		}
	}
	return rolloutStateRollout
}

func (r *ReconcileStatefulSetRollout) update(ctx context.Context, statefulSet *v1beta2.StatefulSet) error {

	err := r.client.Update(ctx, statefulSet)
	if err != nil {
		ctxlog.Errorf(ctx, "Error while updating stateful set: ", err.Error())
		return err
	}
	ctxlog.Debugf(ctx, "StatefulSet %s/%s updated to state Done ", statefulSet.Namespace, statefulSet.Name)
	return nil
}

func readyAndUpdated(ctx context.Context, client crc.Client, statefulSet *v1beta2.StatefulSet) (bool, error) {
	readyAndUpdated := false
	if statefulSet.Spec.UpdateStrategy.RollingUpdate != nil {
		pod, ready, err := GetPod(ctx, client, statefulSet, *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition)
		if err != nil {
			ctxlog.Debug(ctx, "Error calling GetNoneReadyPod ", statefulSet.Namespace, "/", statefulSet.Name, err)
			return false, err
		}
		readyAndUpdated = ready && pod.Labels[v1beta2.StatefulSetRevisionLabel] == statefulSet.Status.UpdateRevision
	}
	return readyAndUpdated, nil
}
