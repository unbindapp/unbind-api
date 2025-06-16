package k8s

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unbindapp/unbind-api/internal/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestContainerStateConstants(t *testing.T) {
	tests := []struct {
		name          string
		state         ContainerState
		expectedValue string
	}{
		{"Running", ContainerStateRunning, "running"},
		{"Waiting", ContainerStateWaiting, "waiting"},
		{"Terminated", ContainerStateTerminated, "terminated"},
		{"Terminating", ContainerStateTerminating, "terminating"},
		{"Crashing", ContainerStateCrashing, "crashing"},
		{"Not Ready", ContainerStateNotReady, "not_ready"},
		{"Image Pull Error", ContainerStateImagePullError, "image_pull_error"},
		{"Starting", ContainerStateStarting, "starting"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, string(tt.state))
		})
	}
}

func TestInstanceHealthConstants(t *testing.T) {
	tests := []struct {
		name          string
		health        InstanceHealth
		expectedValue string
	}{
		{"Pending", InstanceHealthPending, "pending"},
		{"Crashing", InstanceHealthCrashing, "crashing"},
		{"Active", InstanceHealthActive, "active"},
		{"Terminating", InstanceHealthTerminating, "terminating"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, string(tt.health))
		})
	}
}

func TestPodPhaseConstants(t *testing.T) {
	tests := []struct {
		name          string
		phase         PodPhase
		expectedValue string
	}{
		{"Pending", PodPending, "pending"},
		{"Running", PodRunning, "running"},
		{"Succeeded", PodSucceeded, "succeeded"},
		{"Failed", PodFailed, "failed"},
		{"Unknown", PodUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, string(tt.phase))
		})
	}
}

func TestIsPodTerminating(t *testing.T) {
	now := metav1.Now()

	tests := []struct {
		name                string
		pod                 corev1.Pod
		expectedTerminating bool
	}{
		{
			name: "Pod with deletion timestamp",
			pod: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &now,
				},
			},
			expectedTerminating: true,
		},
		{
			name: "Pod in failed state",
			pod: corev1.Pod{
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
				},
			},
			expectedTerminating: false,
		},
		{
			name: "Pod with terminating condition",
			pod: corev1.Pod{
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.PodReady,
							Status: corev1.ConditionFalse,
							Reason: "PodTerminating",
						},
					},
				},
			},
			expectedTerminating: true,
		},
		{
			name: "Normal running pod",
			pod: corev1.Pod{
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.PodReady,
							Status: corev1.ConditionTrue,
							Reason: "PodRunning",
						},
					},
				},
			},
			expectedTerminating: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPodTerminating(tt.pod)
			assert.Equal(t, tt.expectedTerminating, result)
		})
	}
}

func TestMapEventType(t *testing.T) {
	tests := []struct {
		name         string
		reason       string
		message      string
		expectedType models.EventType
	}{
		{
			name:         "OOM Killed event",
			reason:       "OOMKilled",
			message:      "Container was killed due to out of memory",
			expectedType: models.EventTypeOOMKilled,
		},
		{
			name:         "Image pull backoff",
			reason:       "ImagePullBackOff",
			message:      "Unable to pull image",
			expectedType: models.EventTypeImagePullBackOff,
		},
		{
			name:         "Crash loop backoff",
			reason:       "CrashLoopBackOff",
			message:      "Container keeps crashing",
			expectedType: models.EventTypeCrashLoopBackOff,
		},
		{
			name:         "Container created",
			reason:       "Created",
			message:      "Container created successfully",
			expectedType: models.EventTypeContainerCreated,
		},
		{
			name:         "Container started",
			reason:       "Started",
			message:      "Container started successfully",
			expectedType: models.EventTypeContainerStarted,
		},
		{
			name:         "Container stopped",
			reason:       "Killing",
			message:      "Stopping container",
			expectedType: models.EventTypeContainerStopped,
		},
		{
			name:         "Node not ready",
			reason:       "NodeNotReady",
			message:      "Node is not ready",
			expectedType: models.EventTypeNodeNotReady,
		},
		{
			name:         "Failed scheduling",
			reason:       "FailedScheduling",
			message:      "Failed to schedule pod",
			expectedType: models.EventTypeSchedulingFailed,
		},
		{
			name:         "Unknown event",
			reason:       "CustomReason",
			message:      "Custom message",
			expectedType: models.EventTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapEventType(tt.reason, tt.message)
			assert.Equal(t, tt.expectedType, result)
		})
	}
}

func TestMapWaitingReasonToEventType(t *testing.T) {
	tests := []struct {
		name         string
		reason       string
		expectedType models.EventType
	}{
		{
			name:         "Crash loop backoff",
			reason:       "CrashLoopBackOff",
			expectedType: models.EventTypeCrashLoopBackOff,
		},
		{
			name:         "Image pull backoff",
			reason:       "ImagePullBackOff",
			expectedType: models.EventTypeImagePullBackOff,
		},
		{
			name:         "ErrImagePull",
			reason:       "ErrImagePull",
			expectedType: models.EventTypeImagePullBackOff,
		},
		{
			name:         "Container creating",
			reason:       "ContainerCreating",
			expectedType: models.EventTypeContainerCreated,
		},
		{
			name:         "Unknown waiting reason",
			reason:       "CustomWaitingReason",
			expectedType: models.EventTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapWaitingReasonToEventType(tt.reason)
			assert.Equal(t, tt.expectedType, result)
		})
	}
}

func TestFilterEventsByPod(t *testing.T) {
	podName := "test-pod-123"

	events := []models.EventRecord{
		{
			Type:    models.EventTypeContainerStarted,
			Message: "Pod test-pod-123 started successfully",
			Reason:  "Started",
		},
		{
			Type:    models.EventTypeSchedulingFailed,
			Message: "Pod other-pod-456 failed to start",
			Reason:  "Failed",
		},
		{
			Type:    models.EventTypeContainerStarted,
			Message: "Container in test-pod-123 is running",
			Reason:  "ContainerStarted",
		},
		{
			Type:    models.EventTypeUnknown,
			Message: "Unrelated event message",
			Reason:  "Unrelated",
		},
	}

	filtered := filterEventsByPod(events, podName)

	require.Len(t, filtered, 2)
	assert.Contains(t, filtered[0].Message, podName)
	assert.Contains(t, filtered[1].Message, podName)
}

func TestFilterEventsByContainer(t *testing.T) {
	containerName := "app-container"

	events := []models.EventRecord{
		{
			Type:    models.EventTypeContainerStarted,
			Message: "Container app-container started",
			Reason:  "ContainerStarted",
		},
		{
			Type:    models.EventTypeSchedulingFailed,
			Message: "Container other-container failed",
			Reason:  "ContainerFailed",
		},
		{
			Type:    models.EventTypeContainerStarted,
			Message: "app-container is healthy",
			Reason:  "HealthCheck",
		},
		{
			Type:    models.EventTypeUnknown,
			Message: "Network error occurred",
			Reason:  "NetworkError",
		},
	}

	filtered := filterEventsByContainer(events, containerName)

	// The function is permissive and includes:
	// 1. Events mentioning the container name directly
	// 2. Events with "container" in message (generic container events)
	// 3. Events with "started" in reason
	// 4. Events with "failed" or "error" in message
	require.Len(t, filtered, 4) // All events match the permissive criteria

	// Verify events that specifically mention the container name are included
	hasSpecificContainerEvents := 0
	for _, event := range filtered {
		if strings.Contains(event.Message, containerName) {
			hasSpecificContainerEvents++
		}
	}
	assert.Equal(t, 2, hasSpecificContainerEvents) // Two events specifically mention app-container
}

func TestExtractContainerStatus(t *testing.T) {
	podCreatedAt := time.Now().Add(-10 * time.Minute)

	tests := []struct {
		name               string
		containerStatus    corev1.ContainerStatus
		isPodTerminating   bool
		expectedState      ContainerState
		expectedIsCrashing bool
		expectedReady      bool
	}{
		{
			name: "Running and ready container",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        true,
				RestartCount: 0,
				State: corev1.ContainerState{
					Running: &corev1.ContainerStateRunning{
						StartedAt: metav1.Now(),
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateRunning,
			expectedIsCrashing: false,
			expectedReady:      true,
		},
		{
			name: "Running but not ready container",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        false,
				RestartCount: 0,
				State: corev1.ContainerState{
					Running: &corev1.ContainerStateRunning{
						StartedAt: metav1.Now(),
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateNotReady,
			expectedIsCrashing: false,
			expectedReady:      false,
		},
		{
			name: "Container in crash loop backoff",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        false,
				RestartCount: 5,
				State: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{
						Reason:  "CrashLoopBackOff",
						Message: "Container keeps crashing",
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateCrashing,
			expectedIsCrashing: true,
			expectedReady:      false,
		},
		{
			name: "Container with image pull error",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        false,
				RestartCount: 0,
				State: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{
						Reason:  "ImagePullBackOff",
						Message: "Unable to pull image",
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateImagePullError,
			expectedIsCrashing: false,
			expectedReady:      false,
		},
		{
			name: "Container creating",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        false,
				RestartCount: 0,
				State: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{
						Reason:  "ContainerCreating",
						Message: "Creating container",
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateStarting,
			expectedIsCrashing: false,
			expectedReady:      false,
		},
		{
			name: "OOM killed container",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        false,
				RestartCount: 1,
				State: corev1.ContainerState{
					Terminated: &corev1.ContainerStateTerminated{
						ExitCode: 137,
						Reason:   "OOMKilled",
						Message:  "Container was killed due to memory limit",
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateCrashing,
			expectedIsCrashing: true,
			expectedReady:      false,
		},
		{
			name: "Successfully completed container",
			containerStatus: corev1.ContainerStatus{
				Name:         "app-container",
				Ready:        false,
				RestartCount: 0,
				State: corev1.ContainerState{
					Terminated: &corev1.ContainerStateTerminated{
						ExitCode: 0,
						Reason:   "Completed",
						Message:  "Container completed successfully",
					},
				},
			},
			isPodTerminating:   false,
			expectedState:      ContainerStateTerminated,
			expectedIsCrashing: false,
			expectedReady:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContainerStatus(tt.containerStatus, tt.isPodTerminating, podCreatedAt)

			assert.Equal(t, tt.containerStatus.Name, result.KubernetesName)
			assert.Equal(t, tt.expectedReady, result.Ready)
			assert.Equal(t, tt.containerStatus.RestartCount, result.RestartCount)
			assert.Equal(t, tt.expectedState, result.State)
			assert.Equal(t, tt.expectedIsCrashing, result.IsCrashing)
			assert.Equal(t, podCreatedAt, result.PodCreatedAt)
			assert.NotEmpty(t, result.Events)
		})
	}
}

func TestInstanceStatusCreation(t *testing.T) {
	// Test InstanceStatus struct creation and field validation
	status := InstanceStatus{
		KubernetesName:  "test-container",
		Ready:           true,
		RestartCount:    5,
		State:           ContainerStateRunning,
		StateReason:     "Started",
		StateMessage:    "Container is running",
		LastExitCode:    0,
		LastTermination: "Exit code: 0",
		IsCrashing:      false,
		CrashLoopReason: "",
		PodCreatedAt:    time.Now(),
		Events:          []models.EventRecord{},
	}

	assert.Equal(t, "test-container", status.KubernetesName)
	assert.True(t, status.Ready)
	assert.Equal(t, int32(5), status.RestartCount)
	assert.Equal(t, ContainerStateRunning, status.State)
	assert.Equal(t, "Started", status.StateReason)
	assert.Equal(t, "Container is running", status.StateMessage)
	assert.Equal(t, int32(0), status.LastExitCode)
	assert.False(t, status.IsCrashing)
	assert.Empty(t, status.CrashLoopReason)
	assert.NotNil(t, status.Events)
}

func TestPodContainerStatusCreation(t *testing.T) {
	// Test PodContainerStatus struct creation and field validation
	teamID := uuid.New()
	projectID := uuid.New()
	environmentID := uuid.New()
	serviceID := uuid.New()

	status := PodContainerStatus{
		KubernetesName:       "test-pod",
		Namespace:            "default",
		Phase:                PodRunning,
		PodIP:                "10.0.0.1",
		StartTime:            "2023-01-01T00:00:00Z",
		CreatedAt:            time.Now(),
		HasCrashingInstances: false,
		IsTerminating:        false,
		Instances:            []InstanceStatus{},
		InstanceDependencies: []InstanceStatus{},
		TeamID:               teamID,
		ProjectID:            projectID,
		EnvironmentID:        environmentID,
		ServiceID:            serviceID,
	}

	assert.Equal(t, "test-pod", status.KubernetesName)
	assert.Equal(t, "default", status.Namespace)
	assert.Equal(t, PodRunning, status.Phase)
	assert.Equal(t, "10.0.0.1", status.PodIP)
	assert.Equal(t, "2023-01-01T00:00:00Z", status.StartTime)
	assert.False(t, status.HasCrashingInstances)
	assert.False(t, status.IsTerminating)
	assert.Equal(t, teamID, status.TeamID)
	assert.Equal(t, projectID, status.ProjectID)
	assert.Equal(t, environmentID, status.EnvironmentID)
	assert.Equal(t, serviceID, status.ServiceID)
	assert.NotNil(t, status.Instances)
	assert.NotNil(t, status.InstanceDependencies)
}

func TestSimpleHealthStatusCreation(t *testing.T) {
	// Test SimpleHealthStatus struct creation and field validation
	status := SimpleHealthStatus{
		Health:            InstanceHealthActive,
		ExpectedInstances: 3,
		Instances:         []SimpleInstanceStatus{},
	}

	assert.Equal(t, InstanceHealthActive, status.Health)
	assert.Equal(t, 3, status.ExpectedInstances)
	assert.NotNil(t, status.Instances)
	assert.Empty(t, status.Instances)
}

func TestSimpleInstanceStatusCreation(t *testing.T) {
	// Test SimpleInstanceStatus struct creation and field validation
	status := SimpleInstanceStatus{
		KubernetesName: "simple-container",
		Status:         ContainerStateRunning,
		RestartCount:   2,
		PodCreatedAt:   time.Now(),
		Events:         []models.EventRecord{},
	}

	assert.Equal(t, "simple-container", status.KubernetesName)
	assert.Equal(t, ContainerStateRunning, status.Status)
	assert.Equal(t, int32(2), status.RestartCount)
	assert.NotNil(t, status.Events)
}

func TestPodStatusOptions(t *testing.T) {
	// Test struct creation and fields
	options := PodStatusOptions{
		IncludeKubernetesEvents: true,
	}

	assert.True(t, options.IncludeKubernetesEvents)

	options.IncludeKubernetesEvents = false
	assert.False(t, options.IncludeKubernetesEvents)
}

func TestGetPodContainerStatusByLabelsWithFakeClient(t *testing.T) {
	teamID := uuid.New()
	projectID := uuid.New()
	environmentID := uuid.New()
	serviceID := uuid.New()

	pods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "app-pod-1",
				Namespace: "default",
				Labels: map[string]string{
					"app":                "web-server",
					"unbind-team":        teamID.String(),
					"unbind-project":     projectID.String(),
					"unbind-environment": environmentID.String(),
					"unbind-service":     serviceID.String(),
				},
				CreationTimestamp: metav1.Now(),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				PodIP: "10.0.0.1",
				StartTime: &metav1.Time{
					Time: time.Now().Add(-5 * time.Minute),
				},
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "app-container",
						Ready:        true,
						RestartCount: 0,
						State: corev1.ContainerState{
							Running: &corev1.ContainerStateRunning{
								StartedAt: metav1.Now(),
							},
						},
					},
				},
			},
		},
	}

	fakeClient := fake.NewSimpleClientset(pods...)
	kubeClient := &KubeClient{}

	statuses, err := kubeClient.GetPodContainerStatusByLabels(
		context.Background(),
		"default",
		map[string]string{"app": "web-server"},
		fakeClient,
	)
	require.NoError(t, err)
	require.Len(t, statuses, 1)

	pod1 := statuses[0]
	assert.Equal(t, "app-pod-1", pod1.KubernetesName)
	assert.Equal(t, PodRunning, pod1.Phase)
	assert.Equal(t, "10.0.0.1", pod1.PodIP)
	assert.Equal(t, teamID, pod1.TeamID)
	assert.Equal(t, projectID, pod1.ProjectID)
	assert.Equal(t, environmentID, pod1.EnvironmentID)
	assert.Equal(t, serviceID, pod1.ServiceID)
	assert.False(t, pod1.IsTerminating)
	assert.False(t, pod1.HasCrashingInstances)
	assert.Len(t, pod1.Instances, 1)

	container1 := pod1.Instances[0]
	assert.Equal(t, "app-container", container1.KubernetesName)
	assert.True(t, container1.Ready)
	assert.Equal(t, int32(0), container1.RestartCount)
	assert.Equal(t, ContainerStateRunning, container1.State)
	assert.False(t, container1.IsCrashing)
}

func TestGetSimpleHealthStatusWithFakeClient(t *testing.T) {
	pods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "healthy-pod",
				Namespace:         "default",
				Labels:            map[string]string{"app": "test-app"},
				CreationTimestamp: metav1.Now(),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "app-container",
						Ready:        true,
						RestartCount: 0,
						State: corev1.ContainerState{
							Running: &corev1.ContainerStateRunning{
								StartedAt: metav1.Now(),
							},
						},
					},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "crashing-pod",
				Namespace:         "default",
				Labels:            map[string]string{"app": "test-app"},
				CreationTimestamp: metav1.Now(),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "app-container",
						Ready:        false,
						RestartCount: 5,
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{
								Reason:  "CrashLoopBackOff",
								Message: "Container is crashing",
							},
						},
					},
				},
			},
		},
	}

	fakeClient := fake.NewSimpleClientset(pods...)
	kubeClient := &KubeClient{}

	expectedReplicas := 2
	healthStatus, err := kubeClient.GetSimpleHealthStatus(
		context.Background(),
		"default",
		map[string]string{"app": "test-app"},
		&expectedReplicas,
		fakeClient,
	)
	require.NoError(t, err)
	require.NotNil(t, healthStatus)

	assert.Equal(t, 2, healthStatus.ExpectedInstances)
	assert.Len(t, healthStatus.Instances, 2)
	assert.Equal(t, InstanceHealthCrashing, healthStatus.Health)

	var healthyFound, crashingFound bool
	for _, instance := range healthStatus.Instances {
		if instance.KubernetesName == "app-container" && instance.Status == ContainerStateRunning {
			healthyFound = true
		}
		if instance.KubernetesName == "app-container" && instance.Status == ContainerStateCrashing {
			crashingFound = true
		}
	}

	assert.True(t, healthyFound)
	assert.True(t, crashingFound)
}

func TestMapKubernetesPodPhase(t *testing.T) {
	tests := []struct {
		name          string
		kubePhase     corev1.PodPhase
		expectedPhase PodPhase
	}{
		{
			name:          "Kubernetes Pending to our pending",
			kubePhase:     corev1.PodPending,
			expectedPhase: PodPending,
		},
		{
			name:          "Kubernetes Running to our running",
			kubePhase:     corev1.PodRunning,
			expectedPhase: PodRunning,
		},
		{
			name:          "Kubernetes Succeeded to our succeeded",
			kubePhase:     corev1.PodSucceeded,
			expectedPhase: PodSucceeded,
		},
		{
			name:          "Kubernetes Failed to our failed",
			kubePhase:     corev1.PodFailed,
			expectedPhase: PodFailed,
		},
		{
			name:          "Kubernetes Unknown to our unknown",
			kubePhase:     corev1.PodUnknown,
			expectedPhase: PodUnknown,
		},
		{
			name:          "Invalid phase defaults to unknown",
			kubePhase:     corev1.PodPhase("InvalidPhase"),
			expectedPhase: PodUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapKubernetesPodPhase(tt.kubePhase)
			assert.Equal(t, tt.expectedPhase, result)
		})
	}
}
