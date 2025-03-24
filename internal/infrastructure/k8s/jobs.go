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

func (self *KubeClient) CreateDeployment(ctx context.Context, serviceID string, jobID string, env map[string]string) (jobName string, err error) {
	// Cancel any active job for this service
	if err = self.CancelJobsByServiceID(ctx, serviceID); err != nil {
		return "", err
	}

	// Build a unique job name
	jobName = fmt.Sprintf("%s-deployment-%d", jobID, time.Now().Unix())

	// Convert environment variables from map to slice
	var envVars []corev1.EnvVar
	for k, v := range env {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	// Create a default set of annotations
	annotations := map[string]string{
		"deployment-triggered-by": "github-webhook",
		"trigger-timestamp":       time.Now().Format(time.RFC3339),
	}

	// Define the Job object
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"deployment": "true",
				"jobID":      jobID,
				"serviceID":  serviceID,
				"job-name":   jobName,
			},
			Annotations: annotations,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: utils.ToPtr[int32](300),
			BackoffLimit:            utils.ToPtr[int32](0),
			Parallelism:             utils.ToPtr[int32](1),
			// Timeout
			ActiveDeadlineSeconds: utils.ToPtr[int64](1200),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deployment": "true",
						"jobID":      jobID,
						"serviceID":  serviceID,
						"job-name":   jobName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:    "builder-serviceaccount",
					RestartPolicy:         corev1.RestartPolicyNever,
					ShareProcessNamespace: utils.ToPtr(true), // Share process namespace between containers
					Containers: []corev1.Container{
						{
							Name:  "build-container",
							Image: self.config.BuildImage,
							Command: []string{
								"sh",
								"-c",
								`# Wait for buildkitd to be ready
	until buildctl --addr tcp://localhost:1234 debug workers >/dev/null 2>&1; do
		echo "Waiting for BuildKit daemon to be ready...";
		sleep 1;
	done;
	
	# Wait for Docker daemon to be ready
	until docker info >/dev/null 2>&1; do
		echo "Waiting for Docker daemon to be ready...";
		sleep 1;
	done;
	
	# Run the builder with BuildKit
	exec /app/builder`,
							},
							Env: append(envVars, []corev1.EnvVar{
								{
									Name:  "BUILDKIT_HOST",
									Value: "tcp://localhost:1234",
								},
								{
									Name:  "DOCKER_HOST",
									Value: "tcp://localhost:2375",
								},
							}...),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "buildkit-socket",
									MountPath: "/run/buildkit",
								},
							},
						},
						{
							Name:  "buildkit-daemon",
							Image: "moby/buildkit:v0.20.1-rootless",
							Command: []string{
								"sh",
								"-c",
								`
				# Trap SIGTERM/SIGINT and forward it to buildkitd
				trap "echo Received SIGTERM, forwarding to buildkitd; kill -TERM $child" SIGTERM SIGINT
				
				# Start buildkitd in the background
				rootlesskit buildkitd --addr tcp://0.0.0.0:1234 --oci-worker-no-process-sandbox &
				child=$!
				
				# Wait for buildkitd to exit
				wait $child
				status=$?
				
				# If the exit code is 1 (from our SIGTERM), convert it to 0
				if [ $status -eq 1 ]; then
						echo "Buildkitd exited with status 1, converting to 0 to signal graceful shutdown"
						exit 0
				else
						exit $status
				fi
								`,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "BUILDKIT_STEP_LOG_MAX_SIZE",
									Value: "-1", // Disable truncating logs
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: utils.ToPtr(true),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "buildkit-socket",
									MountPath: "/run/buildkit",
								},
								{
									Name:      "buildkit-storage",
									MountPath: "/var/lib/buildkit",
								},
								{
									Name:      "cache-dir",
									MountPath: "/cache",
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"buildctl",
											"--addr", "tcp://localhost:1234",
											"debug", "workers",
										},
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       3,
							},
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
								`# Install necessary tools
	apk add --no-cache procps
	
	# Wait for the /app/builder process to start
	while ! pgrep -f "/app/builder" > /dev/null; do
			echo "Waiting for builder to start..."
			sleep 1
	done
	
	# Get the PID of the actual builder process
	builder_pid=$(pgrep -f "/app/builder")
	echo "Builder process started with PID: $builder_pid"
	
	# Monitor the builder process
	while kill -0 $builder_pid 2>/dev/null; do
			sleep 2
	done
	
	echo "Builder process has completed"
	
	# Give a small grace period for any cleanup
	sleep 5
	
	# Once build is complete, gracefully stop the Docker daemon and buildkit
	echo "Stopping docker and buildkit daemons..."
	pkill -15 dockerd || true
	pkill -15 buildkitd || true
	sleep 3
	
	# Force kill if still running
	pkill -9 dockerd || true 
	pkill -9 buildkitd || true
	echo "Build complete, Docker and BuildKit daemons stopped"
	
	# Exit with success to mark the job as complete
	exit 0`,
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: utils.ToPtr(true),
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "buildkit-socket",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "buildkit-storage",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "docker-config",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "cache-dir",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
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

	// Create the Job in Kubernetes
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
