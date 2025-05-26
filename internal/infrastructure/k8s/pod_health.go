package k8s

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPodContainerStatusByLabels returns pods and their instance status information matching the provided labels
func (self *KubeClient) GetPodContainerStatusByLabels(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) ([]PodContainerStatus, error) {
	// Get pods by labels
	pods, err := self.GetPodsByLabels(ctx, namespace, labels, client)
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

		// Get pod events
		podEvents, err := self.getPodEvents(ctx, namespace, pod.Name, client)
		if err != nil {
			// Log the error but continue processing
			fmt.Printf("Warning: failed to get events for pod %s: %v\n", pod.Name, err)
		}

		// Process regular instances
		hasCrashing := false
		for _, container := range pod.Status.ContainerStatuses {
			instanceStatus := extractContainerStatus(container)
			// Add container-specific events
			instanceStatus.Events = filterEventsByContainer(podEvents, container.Name)

			if instanceStatus.IsCrashing {
				hasCrashing = true
			}
			podStatus.Instances = append(podStatus.Instances, instanceStatus)
		}

		// Process instance dependencies
		for _, container := range pod.Status.InitContainerStatuses {
			instanceStatus := extractContainerStatus(container)
			// Add container-specific events
			instanceStatus.Events = filterEventsByContainer(podEvents, container.Name)

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

// getPodEvents retrieves events for a specific pod
func (self *KubeClient) getPodEvents(ctx context.Context, namespace, podName string, client *kubernetes.Clientset) ([]models.EventRecord, error) {
	// Field selector to filter events for a specific pod
	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s,involvedObject.kind=Pod", podName, namespace)

	// Get events
	events, err := client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	result := make([]models.EventRecord, 0, len(events.Items))

	for _, event := range events.Items {
		eventType := mapEventType(event.Reason, event.Message)

		record := models.EventRecord{
			Type:      eventType,
			Timestamp: event.LastTimestamp.Format(time.RFC3339),
			Message:   event.Message,
			Count:     event.Count,
			Reason:    event.Reason,
		}

		if !event.FirstTimestamp.IsZero() {
			record.FirstSeen = event.FirstTimestamp.Format(time.RFC3339)
		}

		if !event.LastTimestamp.IsZero() {
			record.LastSeen = event.LastTimestamp.Format(time.RFC3339)
		}

		result = append(result, record)
	}

	return result, nil
}

// mapEventType maps Kubernetes event reasons to our EventType enum
func mapEventType(reason, message string) models.EventType {
	reasonLower := strings.ToLower(reason)
	messageLower := strings.ToLower(message)

	switch {
	case strings.Contains(reasonLower, "oom") || strings.Contains(messageLower, "out of memory"):
		return models.EventTypeOOMKilled
	case strings.Contains(reasonLower, "backoff") || strings.Contains(reasonLower, "crashloop"):
		return models.EventTypeCrashLoopBackOff
	case strings.Contains(reasonLower, "created"):
		return models.EventTypeContainerCreated
	case strings.Contains(reasonLower, "started"):
		return models.EventTypeContainerStarted
	case strings.Contains(reasonLower, "killed") || strings.Contains(reasonLower, "stopped"):
		return models.EventTypeContainerStopped
	case strings.Contains(reasonLower, "nodenotready"):
		return models.EventTypeNodeNotReady
	case strings.Contains(reasonLower, "failedscheduling"):
		return models.EventTypeSchedulingFailed
	default:
		return models.EventTypeUnknown
	}
}

// filterEventsByContainer filters pod events to only include those relevant to a specific container
func filterEventsByContainer(events []models.EventRecord, containerName string) []models.EventRecord {
	result := make([]models.EventRecord, 0)

	for _, event := range events {
		// Check if the event message contains the container name
		// This is a heuristic since Kubernetes events don't always clearly indicate which container they belong to
		if strings.Contains(event.Message, containerName) {
			result = append(result, event)
		}
	}

	return result
}

// extractContainerStatus extracts status details from a container status
func extractContainerStatus(container corev1.ContainerStatus) InstanceStatus {
	status := InstanceStatus{
		KubernetesName: container.Name,
		Ready:          container.Ready,
		RestartCount:   container.RestartCount,
		IsCrashing:     false,                  // Default to false, will check crash conditions below
		Events:         []models.EventRecord{}, // Initialize empty events array
	}

	// Determine container state
	switch {
	case container.State.Running != nil:
		status.State = ContainerStateRunning

		// Add event for running container
		if container.State.Running.StartedAt.Time != (time.Time{}) {
			status.Events = append(status.Events, models.EventRecord{
				Type:      models.EventTypeContainerStarted,
				Timestamp: container.State.Running.StartedAt.Format(time.RFC3339),
				Message:   fmt.Sprintf("Container %s started", container.Name),
			})
		}
	case container.State.Waiting != nil:
		status.State = ContainerStateWaiting
		status.StateReason = container.State.Waiting.Reason
		status.StateMessage = container.State.Waiting.Message

		// Check for CrashLoopBackOff condition
		if strings.EqualFold(container.State.Waiting.Reason, "CrashLoopBackOff") {
			status.IsCrashing = true
			status.CrashLoopReason = container.State.Waiting.Message

			// Add CrashLoopBackOff event
			status.Events = append(status.Events, models.EventRecord{
				Type:      models.EventTypeCrashLoopBackOff,
				Timestamp: time.Now().Format(time.RFC3339), // No timestamp in waiting state, use current time
				Message:   container.State.Waiting.Message,
				Reason:    container.State.Waiting.Reason,
			})
		}
	case container.State.Terminated != nil:
		status.State = ContainerStateTerminated
		status.StateReason = container.State.Terminated.Reason
		status.StateMessage = container.State.Terminated.Message
		status.LastExitCode = container.State.Terminated.ExitCode

		// Add terminated event
		terminatedEvent := models.EventRecord{
			Timestamp: container.State.Terminated.FinishedAt.Format(time.RFC3339),
			Message:   container.State.Terminated.Message,
			Reason:    container.State.Terminated.Reason,
		}

		// Check if OOMKilled
		if strings.EqualFold(container.State.Terminated.Reason, "OOMKilled") {
			terminatedEvent.Type = models.EventTypeOOMKilled
		} else {
			terminatedEvent.Type = models.EventTypeContainerStopped
		}

		status.Events = append(status.Events, terminatedEvent)

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

		// Add last termination event
		lastTermEvent := models.EventRecord{
			Timestamp: term.FinishedAt.Format(time.RFC3339),
			Message:   term.Message,
			Reason:    term.Reason,
		}

		// Check if OOMKilled
		if strings.EqualFold(term.Reason, "OOMKilled") {
			lastTermEvent.Type = models.EventTypeOOMKilled
		} else {
			lastTermEvent.Type = models.EventTypeContainerStopped
		}

		status.Events = append(status.Events, lastTermEvent)

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
	KubernetesName  string               `json:"kubernetes_name"`
	Ready           bool                 `json:"ready"`
	RestartCount    int32                `json:"restartCount"`
	State           ContainerState       `json:"state"`
	StateReason     string               `json:"stateReason,omitempty"`
	StateMessage    string               `json:"stateMessage,omitempty"`
	LastExitCode    int32                `json:"lastExitCode,omitempty"`
	LastTermination string               `json:"lastTermination,omitempty"`
	IsCrashing      bool                 `json:"isCrashing"`
	CrashLoopReason string               `json:"crashLoopReason,omitempty"`
	Events          []models.EventRecord `json:"events,omitempty" nullable:"false"`
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
	Health            InstanceHealth         `json:"health"`
	ExpectedInstances int                    `json:"expectedInstances"`
	Instances         []SimpleInstanceStatus `json:"instances"`
}

// SimpleInstanceStatus provides basic instance status information
type SimpleInstanceStatus struct {
	KubernetesName string               `json:"kubernetes_name"`
	Status         ContainerState       `json:"status"` // "Running", "Waiting", "Terminated"
	Events         []models.EventRecord `json:"events,omitempty" nullable:"false"`
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

// Overall health one
type InstanceHealth string

const (
	InstanceHealthHealthy   InstanceHealth = "healthy"
	InstanceHealthDegraded  InstanceHealth = "degraded"
	InstanceHealthUnhealthy InstanceHealth = "unhealthy"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u InstanceHealth) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["InstanceHealth"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "InstanceHealth")
		schemaRef.Title = "InstanceHealth"
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceHealthHealthy))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceHealthDegraded))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceHealthUnhealthy))
		r.Map()["InstanceHealth"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/InstanceHealth"}
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
			Health:            InstanceHealthUnhealthy,
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
	health := InstanceHealthHealthy
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
				Status:         instance.State,
				Events:         instance.Events,
			})
		}
	}

	if hasCrashing {
		health = InstanceHealthUnhealthy
	} else if len(allInstances) < expectedInstances {
		health = InstanceHealthDegraded
	}

	return &SimpleHealthStatus{
		Health:            health,
		ExpectedInstances: expectedInstances,
		Instances:         allInstances,
	}, nil
}
