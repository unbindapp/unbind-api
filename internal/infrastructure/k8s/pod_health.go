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
	"k8s.io/client-go/kubernetes"
)

// GetPodContainerStatusByLabels returns pods and their container status information matching the provided labels
func (k *KubeClient) GetPodContainerStatusByLabels(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) ([]PodContainerStatus, error) {
	// Get pods by labels
	pods, err := k.GetPodsByLabels(ctx, namespace, labels, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	result := make([]PodContainerStatus, 0, len(pods.Items))

	// Extract container status information for each pod
	for _, pod := range pods.Items {
		serviceID, _ := uuid.Parse(pod.Labels["unbind-service"])
		environmentID, _ := uuid.Parse(pod.Labels["unbind-environment"])
		projectID, _ := uuid.Parse(pod.Labels["unbind-project"])
		teamID, _ := uuid.Parse(pod.Labels["unbind-team"])

		podStatus := PodContainerStatus{
			Name:           pod.Name,
			Namespace:      pod.Namespace,
			Phase:          PodPhase(pod.Status.Phase),
			PodIP:          pod.Status.PodIP,
			Containers:     make([]ContainerStatus, 0, len(pod.Status.ContainerStatuses)),
			InitContainers: make([]ContainerStatus, 0, len(pod.Status.InitContainerStatuses)),
			TeamID:         teamID,
			ProjectID:      projectID,
			EnvironmentID:  environmentID,
			ServiceID:      serviceID,
		}

		if pod.Status.StartTime != nil {
			podStatus.StartTime = pod.Status.StartTime.Format(time.RFC3339)
		}

		// Process regular containers
		hasCrashing := false
		for _, container := range pod.Status.ContainerStatuses {
			containerStatus := extractContainerStatus(container)
			if containerStatus.IsCrashing {
				hasCrashing = true
			}
			podStatus.Containers = append(podStatus.Containers, containerStatus)
		}

		// Process init containers
		for _, container := range pod.Status.InitContainerStatuses {
			containerStatus := extractContainerStatus(container)
			if containerStatus.IsCrashing {
				hasCrashing = true
			}
			podStatus.InitContainers = append(podStatus.InitContainers, containerStatus)
		}

		podStatus.HasCrashingContainers = hasCrashing

		result = append(result, podStatus)
	}

	return result, nil
}

// extractContainerStatus extracts status details from a container status
func extractContainerStatus(container corev1.ContainerStatus) ContainerStatus {
	status := ContainerStatus{
		Name:         container.Name,
		Ready:        container.Ready,
		RestartCount: container.RestartCount,
		IsCrashing:   false, // Default to false, will check crash conditions below
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
// ContainerStatus contains information about a container's state
type ContainerStatus struct {
	Name            string         `json:"name"`
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

// PodContainerStatus contains information about a pod and its containers
type PodContainerStatus struct {
	Name                  string            `json:"name"`
	Namespace             string            `json:"namespace"`
	Phase                 PodPhase          `json:"phase"` // Pending, Running, Succeeded, Failed, Unknown
	PodIP                 string            `json:"podIP,omitempty"`
	StartTime             string            `json:"startTime,omitempty"`
	HasCrashingContainers bool              `json:"hasCrashingContainers"`
	Containers            []ContainerStatus `json:"containers" nullable:"false"`
	InitContainers        []ContainerStatus `json:"initContainers" nullable:"false"`
	TeamID                uuid.UUID         `json:"team_id"`
	ProjectID             uuid.UUID         `json:"project_id"`
	EnvironmentID         uuid.UUID         `json:"environment_id"`
	ServiceID             uuid.UUID         `json:"service_id"`
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
