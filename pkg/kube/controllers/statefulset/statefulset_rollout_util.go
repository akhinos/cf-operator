package statefulset

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"code.cloudfoundry.org/quarks-utils/pkg/pointers"

	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"code.cloudfoundry.org/quarks-utils/pkg/ctxlog"
	podutil "code.cloudfoundry.org/quarks-utils/pkg/pod"
	"k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	crc "sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigureStatefulSetForRollout configures a stateful set for canarying and rollout
func ConfigureStatefulSetForRollout(statefulSet *v1beta2.StatefulSet) {
	statefulSet.Spec.UpdateStrategy.Type = v1beta2.RollingUpdateStatefulSetStrategyType
	//the canary rollout is for now directly started, the might move to a webhook instead
	statefulSet.Spec.UpdateStrategy.RollingUpdate = &v1beta2.RollingUpdateStatefulSetStrategy{
		Partition: pointers.Int32(maxInt32(minInt32(*statefulSet.Spec.Replicas, statefulSet.Status.Replicas)-1, 0)),
	}
	statefulSet.Annotations[annotationCanaryRollout] = rolloutStatePending
	statefulSet.Annotations[annotationUpdateStartTime] = strconv.FormatInt(time.Now().Unix(), 10)
}

// FilterLabels filters out labels, that are not suitable for StatefulSet updates
func FilterLabels(labels map[string]string) map[string]string {

	statefulSetLabels := make(map[string]string)
	for key, value := range labels {
		if key != manifest.LabelDeploymentVersion {
			statefulSetLabels[key] = value
		}
	}
	return statefulSetLabels
}

//ComputeAnnotations computes annotations from the instance group
func ComputeAnnotations(ig *manifest.InstanceGroup) (map[string]string, error) {
	statefulSetAnnotations := ig.Env.AgentEnvBoshConfig.Agent.Settings.Annotations
	if statefulSetAnnotations == nil {
		statefulSetAnnotations = make(map[string]string)
	}

	canaryWatchTime, err := ExtractWatchTime(ig.Update.CanaryWatchTime, "canary_watch_time")
	if err != nil {
		return nil, err
	}
	statefulSetAnnotations[annotationCanaryWatchTime] = canaryWatchTime

	updateWatchTime, err := ExtractWatchTime(ig.Update.UpdateWatchTime, "update_watch_time")
	if err != nil {
		return nil, err
	}
	statefulSetAnnotations[annotationUpdateWatchTime] = updateWatchTime

	return statefulSetAnnotations, nil
}

//ExtractWatchTime computes the watch time from a range or an absolute value
func ExtractWatchTime(rawWatchTime string, field string) (string, error) {
	rangeRegex := regexp.MustCompile(`^\s*(\d+)\s*-\s*(\d+)\s*$`)
	if matches := rangeRegex.FindStringSubmatch(rawWatchTime); len(matches) > 0 {
		// Ignore the lower boundary, because the API-Server triggers reconciles
		return matches[2], nil
	}
	absoluteRegex := regexp.MustCompile(`^\s*(\d+)\s*$`)
	if matches := absoluteRegex.FindStringSubmatch(rawWatchTime); len(matches) > 0 {
		return matches[1], nil
	}
	return "", fmt.Errorf("invalid %s", field)
}

// CleanupNonReadyPod deletes all pods, that are not ready
func CleanupNonReadyPod(ctx context.Context, client crc.Client, statefulSet *v1beta2.StatefulSet, index int32) error {
	ctxlog.Debug(ctx, "Cleaning up non ready pod for StatefulSet ", statefulSet.Namespace, "/", statefulSet.Name, "-", index)
	pod, ready, err := GetPod(ctx, client, statefulSet, index)
	if err != nil {
		return err
	}
	if ready || pod == nil {
		return nil
	}
	ctxlog.Debug(ctx, "Deleting pod ", pod.Name)
	return client.Delete(ctx, pod)
}

// GetPod returns a pod for a given statefulset and index
func GetPod(ctx context.Context, client crc.Client, statefulSet *v1beta2.StatefulSet, index int32) (*corev1.Pod, bool, error) {
	var pod corev1.Pod
	podName := fmt.Sprintf("%s-%d", statefulSet.Name, index)
	err := client.Get(ctx, crc.ObjectKey{Name: podName, Namespace: statefulSet.Namespace}, &pod)
	if err != nil {
		if crc.IgnoreNotFound(err) == nil {
			ctxlog.Error(ctx, "Pods ", podName, " belonging to StatefulSet not found", statefulSet.Name, ":", err)
			return nil, false, nil
		}
		return nil, false, err
	}
	return &pod, podutil.IsPodReady(&pod), nil
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
