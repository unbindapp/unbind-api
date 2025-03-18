package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *KubeClient) CreateBuildJob(ctx context.Context, serviceID string, jobID string, env map[string]string) (jobName string, err error) {
	// Cancel any active job for this service.
	if err = self.CancelJobsByServiceID(ctx, serviceID); err != nil {
		return "", err
	}

	// Build a unique job name.
	jobName = fmt.Sprintf("%s-build-%d", jobID, time.Now().Unix())

	// Convert environment variables from map to slice.
	var envVars []corev1.EnvVar
	for k, v := range env {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	// Create a default set of annotations (you can extend these as needed).
	annotations := map[string]string{
		"build-triggered-by": "github-webhook",
		"trigger-timestamp":  time.Now().Format(time.RFC3339),
	}

	// Define the Job object.
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"build":     "true",
				"jobID":     jobID,
				"serviceID": serviceID,
				"job-name":  jobName,
			},
			Annotations: annotations,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: utils.ToPtr[int32](300),
			BackoffLimit:            utils.ToPtr[int32](0),
			Parallelism:             utils.ToPtr[int32](1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"build":     "true",
						"jobID":     jobID,
						"serviceID": serviceID,
						"job-name":  jobName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "builder-serviceaccount",
					RestartPolicy:      corev1.RestartPolicyNever,
					// Share process namespace between containers to allow the monitor to see other containers' processes
					ShareProcessNamespace: utils.ToPtr(true),
					Containers: []corev1.Container{
						{
							Name:  "build-container",
							Image: self.config.BuildImage,
							Command: []string{
								"sh",
								"-c",
								`until docker info > /dev/null 2>&1; do 
    echo "Waiting for Docker daemon to be ready..."; 
    sleep 1; 
done; 
exec /app/builder`,
							},
							Env: append(envVars, corev1.EnvVar{
								Name:  "DOCKER_HOST",
								Value: "tcp://localhost:2375",
							}),
						},
						{
							Name:  "docker-daemon",
							Image: "docker:27.5-dind",
							Env: []corev1.EnvVar{
								{
									Name:  "DOCKER_TLS_CERTDIR",
									Value: "",
								},
							},
							// Args: []string{
							// 	fmt.Sprintf("--insecure-registry=%s", k.config.ContainerRegistryHost),
							// },
							SecurityContext: &corev1.SecurityContext{
								Privileged: utils.ToPtr(true),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "docker-graph-storage",
									MountPath: "/var/lib/docker",
								},
							},
						},
						{
							Name:  "cleanup-monitor",
							Image: "alpine",
							Command: []string{
								"sh",
								"-c",
								`# Wait for the /app/builder process to start
while ! pgrep -f "/app/builder" > /dev/null; do
    sleep 1
done

# Get the PID of the actual builder process
builder_pid=$(pgrep -f "/app/builder")

# Monitor the builder process
while kill -0 $builder_pid 2>/dev/null; do
    sleep 1
done

# Give a small grace period for any cleanup
sleep 5

# Once build is complete, gracefully stop the Docker daemon
pkill dockerd
echo "Build complete, Docker daemon stopped"`,
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: utils.ToPtr(true),
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "docker-graph-storage",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	// Create the Job in Kubernetes.
	_, err = self.clientset.BatchV1().Jobs(self.config.BuilderNamespace).Create(ctx, job, metav1.CreateOptions{})
	return jobName, err
}

// For canceling jobs.
func (self *KubeClient) CancelJobsByServiceID(ctx context.Context, serviceID string) error {
	jobList, err := self.clientset.BatchV1().Jobs(self.config.BuilderNamespace).List(ctx, metav1.ListOptions{
		// We use the "serviceID" label to select jobs.
		LabelSelector: fmt.Sprintf("serviceID=%s", serviceID),
	})
	if err != nil {
		return fmt.Errorf("failed to list jobs for service ID %s: %v", serviceID, err)
	}
	for _, job := range jobList.Items {
		if job.Status.Active > 0 {
			// Delete the job. Using foreground deletion ensures that the pods are cleaned up.
			deletePolicy := metav1.DeletePropagationForeground
			if err := self.clientset.BatchV1().Jobs(self.config.BuilderNamespace).Delete(ctx, job.Name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				return fmt.Errorf("failed to delete job %s: %v", job.Name, err)
			}
			log.Infof("Canceled existing job %s for service %s\n", job.Name, serviceID)
		}
	}
	return nil
}

// Get status of a kubernetes Job resource
type JobConditionType int

const (
	JobComplete JobConditionType = iota
	JobFailed
	JobRunning
	JobPending
)

func (js JobConditionType) String() string {
	return [...]string{"JobComplete", "JobFailed", "JobRunning", "JobPending"}[js]
}

// JobStatus represents the status of a job with additional details
type JobStatus struct {
	ConditionType JobConditionType
	FailureReason string
}

func (self *KubeClient) GetJobStatus(ctx context.Context, jobName string) (JobStatus, error) {
	// Get the job from Kubernetes API
	job, err := self.clientset.BatchV1().Jobs(self.config.BuilderNamespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return JobStatus{}, fmt.Errorf("failed to get job %s: %v", jobName, err)
	}

	result := JobStatus{}

	// Determine status based on job conditions
	if job.Status.Succeeded > 0 {
		result.ConditionType = JobComplete
	} else if job.Status.Failed > 0 {
		result.ConditionType = JobFailed
		// Try to extract failure reason from conditions
		for _, condition := range job.Status.Conditions {
			if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
				result.FailureReason = condition.Reason + ": " + condition.Message
				break
			}
		}

		// If no reason found in conditions, try to get it from the associated pods
		if result.FailureReason == "" {
			result.FailureReason = self.getJobPodsFailureReason(ctx, job.Name)
		}
	} else if job.Status.Active > 0 {
		result.ConditionType = JobRunning
	} else {
		result.ConditionType = JobPending
	}

	return result, nil
}

// getJobPodsFailureReason attempts to get failure reasons from pods associated with the job
func (self *KubeClient) getJobPodsFailureReason(ctx context.Context, jobName string) string {
	// Get pods with the job-name label
	labelSelector := fmt.Sprintf("job-name=%s", jobName)
	pods, err := self.clientset.CoreV1().Pods(self.config.BuilderNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil || len(pods.Items) == 0 {
		return "Unknown failure reason"
	}

	// Look for failed pods and extract their termination reason
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.ExitCode != 0 {
				reason := containerStatus.State.Terminated.Reason
				message := containerStatus.State.Terminated.Message

				if message != "" {
					return fmt.Sprintf("%s: %s (Container: %s, Exit Code: %d)",
						reason, message, containerStatus.Name, containerStatus.State.Terminated.ExitCode)
				}

				return fmt.Sprintf("%s (Container: %s, Exit Code: %d)",
					reason, containerStatus.Name, containerStatus.State.Terminated.ExitCode)
			}
		}

		// If no terminated containers found, check pod's status
		if pod.Status.Phase == corev1.PodFailed {
			return fmt.Sprintf("Pod failed: %s", pod.Status.Reason)
		}
	}

	return "Unknown failure reason"
}
