package buildctl

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/infrastructure/queue"
)

// The request to build a service, includes environment for builder image
type BuildJobRequest struct {
	ServiceID   uuid.UUID         `json:"service_id"`
	Environment map[string]string `json:"environment"`
}

// Handles triggering builds for services
type BuildController struct {
	k8s      *k8s.KubeClient
	jobQueue *queue.Queue[BuildJobRequest]
}
