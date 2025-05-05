package k8s

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPodContainerStatusByLabels returns pods and their instance status information matching the provided labels
func (k *KubeClient) GetPodContainerStatusByLabels(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) ([]PodContainerStatus, error) {
	// Get pods by labels
	pods, err := k.GetPodsByLabels(ctx, namespace, labels, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	result := make([]PodContainerStatus, 0, len(pods.Items))

	// Extract instance status information for each pod
	for _, pod := range pods.Items {
		serviceID, _ := uuid.Parse(pod.Labels["unbind-service"])
		environmentID, _ := uuid.Parse(pod.Labels["unbind-environment"])
		projectID, _ := uuid.Parse(pod.Labels["unbind-project"])
		teamID, _ := uuid.Parse(pod.Labels["unbind-team"])

		podStatus := PodContainerStatus{
			KubernetesName:       pod.Name,
			Namespace:            pod.Namespace,
			Phase:                PodPhase(pod.Status.Phase),
			PodIP:                pod.Status.PodIP,
			Instances:            make([]InstanceStatus, 0, len(pod.Status.ContainerStatuses)),
			InstanceDependencies: make([]InstanceStatus, 0, len(pod.Status.InitContainerStatuses)),
			TeamID:               teamID,
			ProjectID:            projectID,
			EnvironmentID:        environmentID,
			ServiceID:            serviceID,
		}

		if pod.Status.StartTime != nil {
			podStatus.StartTime = pod.Status.StartTime.Format(time.RFC3339)
		}

		// Process regular instances
		hasCrashing := false
		for _, container := range pod.Status.ContainerStatuses {
			instanceStatus := extractContainerStatus(container)
			if instanceStatus.IsCrashing {
				hasCrashing = true
			}
			podStatus.Instances = append(podStatus.Instances, instanceStatus)
		}

		// Process instance dependencies
		for _, container := range pod.Status.InitContainerStatuses {
			instanceStatus := extractContainerStatus(container)
			if instanceStatus.IsCrashing {
				hasCrashing = true
			}
			podStatus.InstanceDependencies = append(podStatus.InstanceDependencies, instanceStatus)
		}

		podStatus.HasCrashingInstances = hasCrashing

		result = append(result, podStatus)
	}

	return result, nil
}

// extractContainerStatus extracts status details from a container status
func extractContainerStatus(container corev1.ContainerStatus) InstanceStatus {
	status := InstanceStatus{
		KubernetesName: container.Name,
		Ready:          container.Ready,
		RestartCount:   container.RestartCount,
		IsCrashing:     false, // Default to false, will check crash conditions below
	}

	// Determine container state
	switch {
	case container.State.Running != nil:
		status.State = ContainerStateRunning
	case container.State.Waiting != nil:
		status.State = ContainerStateWaiting
		status.StateReason = container.State.Waiting.Reason
		status.StateMessage = container.State.Waiting.Message

		// Check for CrashLoopBackOff condition
		if strings.EqualFold(container.State.Waiting.Reason, "CrashLoopBackOff") {
			status.IsCrashing = true
			status.CrashLoopReason = container.State.Waiting.Message
		}
	case container.State.Terminated != nil:
		status.State = ContainerStateTerminated
		status.StateReason = container.State.Terminated.Reason
		status.StateMessage = container.State.Terminated.Message
		status.LastExitCode = container.State.Terminated.ExitCode

		// Check if container terminated with non-zero exit code
		if container.State.Terminated.ExitCode != 0 {
			status.IsCrashing = true
			status.CrashLoopReason = fmt.Sprintf("Terminated with exit code: %d, reason: %s",
				container.State.Terminated.ExitCode, container.State.Terminated.Reason)
		}
	}

	// Get information about last termination if available
	if container.LastTerminationState.Terminated != nil {
		term := container.LastTerminationState.Terminated
		status.LastTermination = fmt.Sprintf(
			"Exit code: %d, Reason: %s, Message: %s, Finished at: %s",
			term.ExitCode,
			term.Reason,
			term.Message,
			term.FinishedAt.Format(time.RFC3339),
		)

		// High restart count and non-zero exit code
		if container.RestartCount > 3 && term.ExitCode != 0 {
			status.IsCrashing = true
			if status.CrashLoopReason == "" {
				status.CrashLoopReason = fmt.Sprintf("Frequent restarts (%d) with exit code: %d",
					container.RestartCount, term.ExitCode)
			}
		}
	}

	return status
}

// * Models
// InstanceStatus contains information about an instance's state
type InstanceStatus struct {
	KubernetesName  string         `json:"kubernetes_name"`
	Ready           bool           `json:"ready"`
	RestartCount    int32          `json:"restartCount"`
	State           ContainerState `json:"state"`
	StateReason     string         `json:"stateReason,omitempty"`
	StateMessage    string         `json:"stateMessage,omitempty"`
	LastExitCode    int32          `json:"lastExitCode,omitempty"`
	LastTermination string         `json:"lastTermination,omitempty"`
	IsCrashing      bool           `json:"isCrashing"`
	CrashLoopReason string         `json:"crashLoopReason,omitempty"`
}

// PodContainerStatus contains information about a pod and its instances
type PodContainerStatus struct {
	KubernetesName       string           `json:"kubernetes_name"`
	Namespace            string           `json:"namespace"`
	Phase                PodPhase         `json:"phase"` // Pending, Running, Succeeded, Failed, Unknown
	PodIP                string           `json:"podIP,omitempty"`
	StartTime            string           `json:"startTime,omitempty"`
	HasCrashingInstances bool             `json:"hasCrashingInstances"`
	Instances            []InstanceStatus `json:"instances" nullable:"false"`
	InstanceDependencies []InstanceStatus `json:"instanceDependencies" nullable:"false"`
	TeamID               uuid.UUID        `json:"team_id"`
	ProjectID            uuid.UUID        `json:"project_id"`
	EnvironmentID        uuid.UUID        `json:"environment_id"`
	ServiceID            uuid.UUID        `json:"service_id"`
}

// SimpleHealthStatus provides a simplified view of instance health
type SimpleHealthStatus struct {
	Health            string                 `json:"health"` // "healthy", "unhealthy", "degraded"
	ExpectedInstances int                    `json:"expectedInstances"`
	Instances         []SimpleInstanceStatus `json:"instances"`
}

// SimpleInstanceStatus provides basic instance status information
type SimpleInstanceStatus struct {
	KubernetesName string `json:"kubernetes_name"`
	Status         string `json:"status"` // "Running", "Waiting", "Terminated"
}

// * Enums
type ContainerState string

const (
	ContainerStateRunning    ContainerState = "Running"
	ContainerStateWaiting    ContainerState = "Waiting"
	ContainerStateTerminated ContainerState = "Terminated"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u ContainerState) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["ContainerState"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "ContainerState")
		schemaRef.Title = "ContainerState"
		schemaRef.Enum = append(schemaRef.Enum, string(ContainerStateRunning))
		schemaRef.Enum = append(schemaRef.Enum, string(ContainerStateWaiting))
		schemaRef.Enum = append(schemaRef.Enum, string(ContainerStateTerminated))
		r.Map()["ContainerState"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/ContainerState"}
}

// Direct copy of corev1.PodPhase, but so we can attach openapi schema to it
type PodPhase string

// These are the valid statuses of pods.
const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
	PodUnknown   PodPhase = "Unknown"
)

var allPodPhases = []corev1.PodPhase{
	corev1.PodPending,
	corev1.PodRunning,
	corev1.PodSucceeded,
	corev1.PodFailed,
	corev1.PodUnknown,
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u PodPhase) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["PodPhase"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "PodPhase")
		schemaRef.Title = "PodPhase"
		for _, v := range allPodPhases {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["PodPhase"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/PodPhase"}
}

// GetExpectedInstances determines the expected number of instances based on the parent resource
func (k *KubeClient) GetExpectedInstances(ctx context.Context, namespace string, podName string, client *kubernetes.Clientset) (int, error) {
	// Get the pod to find its owner reference
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get pod: %w", err)
	}

	// Look for owner references
	for _, ownerRef := range pod.OwnerReferences {
		switch ownerRef.Kind {
		case "StatefulSet":
			sts, err := client.AppsV1().StatefulSets(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get statefulset: %w", err)
			}
			return int(*sts.Spec.Replicas), nil

		case "Deployment":
			deploy, err := client.AppsV1().Deployments(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get deployment: %w", err)
			}
			return int(*deploy.Spec.Replicas), nil

		case "ReplicaSet":
			rs, err := client.AppsV1().ReplicaSets(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get replicaset: %w", err)
			}
			return int(*rs.Spec.Replicas), nil

		case "DaemonSet":
			// For DaemonSets, the expected count is the number of nodes that match the node selector
			ds, err := client.AppsV1().DaemonSets(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get daemonset: %w", err)
			}
			return int(ds.Status.DesiredNumberScheduled), nil
		}
	}

	// If no owner reference found or it's a standalone pod, return 1
	return 1, nil
}

// GetSimpleHealthStatus gets a simplified health status for all pods matching the labels
func (self *KubeClient) GetSimpleHealthStatus(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) (*SimpleHealthStatus, error) {
	// Get all pods matching the labels
	podStatuses, err := self.GetPodContainerStatusByLabels(ctx, namespace, labels, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod statuses: %w", err)
	}

	if len(podStatuses) == 0 {
		return &SimpleHealthStatus{
			Health:            "unhealthy",
			ExpectedInstances: 0,
			Instances:         []SimpleInstanceStatus{},
		}, nil
	}

	// Get expected instances from the first pod (they should all have the same parent)
	expectedInstances, err := self.GetExpectedInstances(ctx, podStatuses[0].Namespace, podStatuses[0].KubernetesName, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get expected instances: %w", err)
	}

	// Determine overall health
	health := "healthy"
	hasCrashing := false
	allInstances := make([]SimpleInstanceStatus, 0)

	for _, podStatus := range podStatuses {
		if podStatus.HasCrashingInstances {
			hasCrashing = true
		}

		// Add all instances from this pod
		for _, instance := range podStatus.Instances {
			allInstances = append(allInstances, SimpleInstanceStatus{
				KubernetesName: instance.KubernetesName,
				Status:         string(instance.State),
			})
		}
	}

	if hasCrashing {
		health = "unhealthy"
	} else if len(allInstances) < expectedInstances {
		health = "degraded"
	}

	return &SimpleHealthStatus{
		Health:            health,
		ExpectedInstances: expectedInstances,
		Instances:         allInstances,
	}, nil
}
