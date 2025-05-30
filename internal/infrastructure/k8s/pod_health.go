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

// GetPodContainerStatusByLabels efficiently fetches pod status with inferred events from container state
func (self *KubeClient) GetPodContainerStatusByLabels(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) ([]PodContainerStatus, error) {
	return self.GetPodContainerStatusByLabelsWithOptions(ctx, namespace, labels, client, PodStatusOptions{
		IncludeKubernetesEvents: false, // Only inferred events by default
	})
}

// PodStatusOptions controls what data to fetch for pod status
type PodStatusOptions struct {
	IncludeKubernetesEvents bool // Whether to fetch additional events from Kubernetes Events API (more expensive)
}

// GetPodContainerStatusByLabelsWithOptions efficiently fetches pod status with configurable options
// Container state events are always inferred (lightweight and reliable)
func (self *KubeClient) GetPodContainerStatusByLabelsWithOptions(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset, options PodStatusOptions) ([]PodContainerStatus, error) {
	pods, err := self.GetPodsByLabels(ctx, namespace, labels, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	result := make([]PodContainerStatus, 0, len(pods.Items))

	// Batch fetch Kubernetes Events API events for all pods if requested
	var allEvents []models.EventRecord
	if options.IncludeKubernetesEvents && len(pods.Items) > 0 {
		allEvents, err = self.getBatchPodEvents(ctx, namespace, pods.Items, client)
		if err != nil {
			// Log warning but continue without Kubernetes events
			fmt.Printf("Warning: failed to get Kubernetes events for pods: %v\n", err)
			allEvents = []models.EventRecord{}
		}
	}

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

		// Filter Kubernetes events for this specific pod if requested
		var podEvents []models.EventRecord
		if options.IncludeKubernetesEvents {
			podEvents = filterEventsByPod(allEvents, pod.Name)
		}

		hasCrashing := false
		for _, container := range pod.Status.ContainerStatuses {
			// Always extract inferred events from container state (lightweight and reliable)
			instanceStatus := extractContainerStatus(container)

			// Optionally append additional Kubernetes Events API events
			if options.IncludeKubernetesEvents {
				filteredEvents := filterEventsByContainer(podEvents, container.Name)
				instanceStatus.Events = append(instanceStatus.Events, filteredEvents...)
			}

			if instanceStatus.IsCrashing {
				hasCrashing = true
			}
			podStatus.Instances = append(podStatus.Instances, instanceStatus)
		}

		for _, container := range pod.Status.InitContainerStatuses {
			// Always extract inferred events from container state
			instanceStatus := extractContainerStatus(container)

			// Optionally append additional Kubernetes Events API events
			if options.IncludeKubernetesEvents {
				filteredEvents := filterEventsByContainer(podEvents, container.Name)
				instanceStatus.Events = append(instanceStatus.Events, filteredEvents...)
			}

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

// getBatchPodEvents efficiently fetches events for multiple pods in a single API call
func (self *KubeClient) getBatchPodEvents(ctx context.Context, namespace string, pods []corev1.Pod, client *kubernetes.Clientset) ([]models.EventRecord, error) {
	result := make([]models.EventRecord, 0)

	// Single API call to get all events in the namespace
	eventsV1, err := client.EventsV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	// Create a map of pod names for efficient lookup
	podNames := make(map[string]bool, len(pods))
	for _, pod := range pods {
		podNames[pod.Name] = true
	}

	// Filter events to only include those related to our pods
	for _, event := range eventsV1.Items {
		isRelevant := false

		// Check if event is directly related to one of our pods
		if strings.EqualFold(event.Regarding.Kind, "Pod") && podNames[event.Regarding.Name] {
			isRelevant = true
		} else {
			// Check if event mentions any of our pod names
			for podName := range podNames {
				if strings.Contains(event.Note, podName) {
					isRelevant = true
					break
				}
				if strings.EqualFold(event.Regarding.Kind, "Container") && strings.HasPrefix(event.Regarding.Name, podName) {
					isRelevant = true
					break
				}
			}
		}

		if isRelevant {
			eventType := mapEventType(event.Reason, event.Note)

			record := models.EventRecord{
				Type:      eventType,
				Timestamp: event.EventTime.Format(time.RFC3339),
				Message:   event.Note,
				Count:     1,
				Reason:    event.Reason,
			}
			if event.Series != nil {
				record.Count = event.Series.Count
			}

			if !event.EventTime.IsZero() {
				record.FirstSeen = event.EventTime.Format(time.RFC3339)
				record.LastSeen = event.EventTime.Format(time.RFC3339)
			}

			result = append(result, record)
		}
	}

	return result, nil
}

// filterEventsByPod filters events to only include those related to a specific pod
func filterEventsByPod(events []models.EventRecord, podName string) []models.EventRecord {
	result := make([]models.EventRecord, 0)

	for _, event := range events {
		if strings.Contains(event.Message, podName) || strings.Contains(event.Reason, podName) {
			result = append(result, event)
		}
	}

	return result
}

func mapEventType(reason, message string) models.EventType {
	reasonLower := strings.ToLower(reason)
	messageLower := strings.ToLower(message)

	switch {
	case strings.Contains(reasonLower, "oom") || strings.Contains(messageLower, "out of memory"):
		return models.EventTypeOOMKilled
	case strings.Contains(reasonLower, "imagepullbackoff"):
		return models.EventTypeImagePullBackOff
	case strings.Contains(reasonLower, "backoff") || strings.Contains(reasonLower, "crashloop"):
		return models.EventTypeCrashLoopBackOff
	case strings.Contains(reasonLower, "created") || reasonLower == "created":
		return models.EventTypeContainerCreated
	case strings.Contains(reasonLower, "started") || reasonLower == "started":
		return models.EventTypeContainerStarted
	case strings.Contains(reasonLower, "pulled") || strings.Contains(messageLower, "successfully pulled image"):
		return models.EventTypeContainerCreated
	case strings.Contains(reasonLower, "killing") || strings.Contains(reasonLower, "killed") || strings.Contains(reasonLower, "stopped"):
		return models.EventTypeContainerStopped
	case strings.Contains(reasonLower, "nodenotready"):
		return models.EventTypeNodeNotReady
	case strings.Contains(reasonLower, "failedscheduling"):
		return models.EventTypeSchedulingFailed
	case strings.Contains(reasonLower, "failed") || strings.Contains(messageLower, "failed"):
		return models.EventTypeUnknown
	default:
		return models.EventTypeUnknown
	}
}

func filterEventsByContainer(events []models.EventRecord, containerName string) []models.EventRecord {
	result := make([]models.EventRecord, 0)

	for _, event := range events {
		if strings.Contains(event.Message, containerName) {
			result = append(result, event)
			continue
		}

		reasonLower := strings.ToLower(event.Reason)
		messageLower := strings.ToLower(event.Message)

		if strings.Contains(reasonLower, "created") ||
			strings.Contains(reasonLower, "started") ||
			strings.Contains(reasonLower, "pulled") ||
			strings.Contains(reasonLower, "killing") ||
			strings.Contains(reasonLower, "preempting") ||
			strings.Contains(messageLower, "container") ||
			strings.Contains(messageLower, "image") {
			result = append(result, event)
			continue
		}

		if event.Type == models.EventTypeOOMKilled ||
			event.Type == models.EventTypeCrashLoopBackOff ||
			event.Type == models.EventTypeContainerStopped ||
			strings.Contains(messageLower, "failed") ||
			strings.Contains(messageLower, "error") ||
			strings.Contains(messageLower, "warning") {
			result = append(result, event)
		}
	}

	return result
}

// extractContainerStatus infers some events from the container status, since Events API may not always have old events
func extractContainerStatus(container corev1.ContainerStatus) InstanceStatus {
	status := InstanceStatus{
		KubernetesName: container.Name,
		Ready:          container.Ready,
		RestartCount:   container.RestartCount,
		IsCrashing:     false,
		Events:         []models.EventRecord{},
	}

	switch {
	case container.State.Running != nil:
		// Container is running, but check if it's ready
		if container.Ready {
			status.State = ContainerStateRunning
		} else {
			// Container is running but not ready (likely failing readiness probes)
			status.State = ContainerStateWaiting
			status.StateReason = "NotReady"
			status.StateMessage = "Container is running but not ready (readiness probe may be failing)"
		}

		if container.State.Running.StartedAt.Time != (time.Time{}) {
			if container.Ready {
				status.Events = append(status.Events, models.EventRecord{
					Type:      models.EventTypeContainerStarted,
					Timestamp: container.State.Running.StartedAt.Format(time.RFC3339),
					Message:   fmt.Sprintf("Container %s started at %s", container.Name, container.State.Running.StartedAt.Format(time.RFC3339)),
					Reason:    "Started",
				})
			} else {
				status.Events = append(status.Events, models.EventRecord{
					Type:      models.EventTypeUnknown,
					Timestamp: container.State.Running.StartedAt.Format(time.RFC3339),
					Message:   fmt.Sprintf("Container %s started but not ready (readiness probe failing)", container.Name),
					Reason:    "NotReady",
				})
			}
		}

		if container.RestartCount > 0 {
			status.Events = append(status.Events, models.EventRecord{
				Type:      models.EventTypeContainerCreated,
				Timestamp: container.State.Running.StartedAt.Format(time.RFC3339),
				Message:   fmt.Sprintf("Container %s restarted (restart count: %d)", container.Name, container.RestartCount),
				Reason:    "Created",
			})
		} else {
			status.Events = append(status.Events, models.EventRecord{
				Type:      models.EventTypeContainerCreated,
				Timestamp: container.State.Running.StartedAt.Format(time.RFC3339),
				Message:   fmt.Sprintf("Container %s created and started", container.Name),
				Reason:    "Created",
			})
		}

	case container.State.Waiting != nil:
		status.State = ContainerStateWaiting
		status.StateReason = container.State.Waiting.Reason
		status.StateMessage = container.State.Waiting.Message

		status.Events = append(status.Events, models.EventRecord{
			Type:      models.EventTypeUnknown,
			Timestamp: time.Now().Format(time.RFC3339),
			Message:   fmt.Sprintf("Container %s is waiting: %s", container.Name, container.State.Waiting.Message),
			Reason:    container.State.Waiting.Reason,
		})

		if strings.EqualFold(container.State.Waiting.Reason, "CrashLoopBackOff") {
			status.IsCrashing = true
			status.CrashLoopReason = container.State.Waiting.Message

			if len(status.Events) > 0 {
				status.Events[len(status.Events)-1].Type = models.EventTypeCrashLoopBackOff
			}
		}

	case container.State.Terminated != nil:
		status.State = ContainerStateTerminated
		status.StateReason = container.State.Terminated.Reason
		status.StateMessage = container.State.Terminated.Message
		status.LastExitCode = container.State.Terminated.ExitCode

		terminatedEvent := models.EventRecord{
			Timestamp: container.State.Terminated.FinishedAt.Format(time.RFC3339),
			Message:   fmt.Sprintf("Container %s terminated with exit code %d: %s", container.Name, container.State.Terminated.ExitCode, container.State.Terminated.Message),
			Reason:    container.State.Terminated.Reason,
		}

		if strings.EqualFold(container.State.Terminated.Reason, "OOMKilled") {
			terminatedEvent.Type = models.EventTypeOOMKilled
		} else {
			terminatedEvent.Type = models.EventTypeContainerStopped
		}

		status.Events = append(status.Events, terminatedEvent)

		if container.State.Terminated.ExitCode != 0 && !strings.EqualFold(container.State.Terminated.Reason, "Completed") {
			status.IsCrashing = true
			status.CrashLoopReason = fmt.Sprintf("Terminated with exit code: %d, reason: %s",
				container.State.Terminated.ExitCode, container.State.Terminated.Reason)
		}
	}

	if container.LastTerminationState.Terminated != nil {
		term := container.LastTerminationState.Terminated
		status.LastTermination = fmt.Sprintf(
			"Exit code: %d, Reason: %s, Message: %s, Finished at: %s",
			term.ExitCode,
			term.Reason,
			term.Message,
			term.FinishedAt.Format(time.RFC3339),
		)

		lastTermEvent := models.EventRecord{
			Timestamp: term.FinishedAt.Format(time.RFC3339),
			Message:   fmt.Sprintf("Container %s previous termination: exit code %d, %s", container.Name, term.ExitCode, term.Message),
			Reason:    term.Reason,
		}

		if strings.EqualFold(term.Reason, "OOMKilled") {
			lastTermEvent.Type = models.EventTypeOOMKilled
		} else {
			lastTermEvent.Type = models.EventTypeContainerStopped
		}

		status.Events = append(status.Events, lastTermEvent)
	}

	return status
}

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

type PodContainerStatus struct {
	KubernetesName       string           `json:"kubernetes_name"`
	Namespace            string           `json:"namespace"`
	Phase                PodPhase         `json:"phase"`
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

type SimpleHealthStatus struct {
	Health            InstanceHealth         `json:"health"`
	ExpectedInstances int                    `json:"expectedInstances"`
	Instances         []SimpleInstanceStatus `json:"instances"`
}

type SimpleInstanceStatus struct {
	KubernetesName string               `json:"kubernetes_name"`
	Status         ContainerState       `json:"status"`
	Events         []models.EventRecord `json:"events,omitempty" nullable:"false"`
}

type ContainerState string

const (
	ContainerStateRunning    ContainerState = "Running"
	ContainerStateWaiting    ContainerState = "Waiting"
	ContainerStateTerminated ContainerState = "Terminated"
)

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

type InstanceHealth string

const (
	InstanceHealthPending  InstanceHealth = "pending"  // Waiting to be scheduled, or running but not ready yet
	InstanceHealthCrashing InstanceHealth = "crashing" // Has crashing instances
	InstanceHealthActive   InstanceHealth = "active"   // All instances running and healthy
)

func (u InstanceHealth) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["InstanceHealth"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "InstanceHealth")
		schemaRef.Title = "InstanceHealth"
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceHealthPending))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceHealthCrashing))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceHealthActive))
		r.Map()["InstanceHealth"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/InstanceHealth"}
}

type PodPhase string

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

func (k *KubeClient) GetExpectedInstances(ctx context.Context, namespace string, podName string, client *kubernetes.Clientset) (int, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get pod: %w", err)
	}

	for _, ownerRef := range pod.OwnerReferences {
		switch strings.ToLower(ownerRef.Kind) {
		case "statefulset":
			sts, err := client.AppsV1().StatefulSets(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get statefulset: %w", err)
			}
			return int(*sts.Spec.Replicas), nil

		case "deployment":
			deploy, err := client.AppsV1().Deployments(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get deployment: %w", err)
			}
			return int(*deploy.Spec.Replicas), nil

		case "replicaset":
			rs, err := client.AppsV1().ReplicaSets(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get replicaset: %w", err)
			}
			return int(*rs.Spec.Replicas), nil

		case "daemonset":
			ds, err := client.AppsV1().DaemonSets(namespace).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to get daemonset: %w", err)
			}
			return int(ds.Status.DesiredNumberScheduled), nil
		}
	}

	return 1, nil
}

func (self *KubeClient) GetSimpleHealthStatus(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) (*SimpleHealthStatus, error) {
	podStatuses, err := self.GetPodContainerStatusByLabels(ctx, namespace, labels, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod statuses: %w", err)
	}

	if len(podStatuses) == 0 {
		return &SimpleHealthStatus{
			Health:            InstanceHealthPending,
			ExpectedInstances: 0,
			Instances:         []SimpleInstanceStatus{},
		}, nil
	}

	expectedInstances, err := self.GetExpectedInstances(ctx, podStatuses[0].Namespace, podStatuses[0].KubernetesName, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get expected instances: %w", err)
	}

	hasCrashing := false
	hasPending := false
	readyCount := 0
	allInstances := make([]SimpleInstanceStatus, 0)

	for _, podStatus := range podStatuses {
		// Check if any containers are crashing
		if podStatus.HasCrashingInstances {
			hasCrashing = true
		}

		for _, instance := range podStatus.Instances {
			allInstances = append(allInstances, SimpleInstanceStatus{
				KubernetesName: instance.KubernetesName,
				Status:         instance.State,
				Events:         instance.Events,
			})

			// Count ready instances and detect pending states
			if instance.IsCrashing {
				hasCrashing = true
			} else if instance.State == ContainerStateRunning && instance.Ready {
				readyCount++
			} else {
				// Container is waiting, terminated, or running but not ready
				hasPending = true
			}
		}
	}

	// Determine health status based on priority:
	// 1. Crashing takes precedence over everything
	// 2. Pending if any containers are not ready or we don't have enough instances
	// 3. Active only if all expected instances are ready
	var health InstanceHealth
	if hasCrashing {
		health = InstanceHealthCrashing
	} else if hasPending || readyCount < expectedInstances {
		health = InstanceHealthPending
	} else {
		health = InstanceHealthActive
	}

	return &SimpleHealthStatus{
		Health:            health,
		ExpectedInstances: expectedInstances,
		Instances:         allInstances,
	}, nil
}
