package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type DeploymentResponse struct {
	ID            uuid.UUID               `json:"id"`
	ServiceID     uuid.UUID               `json:"service_id"`
	Status        schema.DeploymentStatus `json:"status"`
	Error         string                  `json:"error,omitempty"`
	StartedAt     *time.Time              `json:"started_at,omitempty"`
	CompletedAt   *time.Time              `json:"completed_at,omitempty"`
	Attempts      int                     `json:"attempts"`
	CommitSHA     *string                 `json:"commit_sha,omitempty" required:"false"`
	CommitMessage *string                 `json:"commit_message,omitempty" required:"false"`
	CommitAuthor  *schema.GitCommitter    `json:"commit_author,omitempty" required:"false"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// TransformDeploymentEntity transforms an ent.Deployment entity into a DeploymentResponse
func TransformDeploymentEntity(entity *ent.Deployment) *DeploymentResponse {
	response := &DeploymentResponse{}
	if entity != nil {
		response = &DeploymentResponse{
			ID:            entity.ID,
			ServiceID:     entity.ServiceID,
			Status:        entity.Status,
			Error:         entity.Error,
			StartedAt:     entity.StartedAt,
			CompletedAt:   entity.CompletedAt,
			Attempts:      entity.Attempts,
			CommitSHA:     entity.CommitSha,
			CommitMessage: entity.CommitMessage,
			CommitAuthor:  entity.CommitAuthor,
			CreatedAt:     entity.CreatedAt,
			UpdatedAt:     entity.UpdatedAt,
		}
	}
	return response
}

// Transforms a slice of ent.Deployment entities into a slice of DeploymentResponse
func TransformDeploymentEntities(entities []*ent.Deployment) []*DeploymentResponse {
	responses := make([]*DeploymentResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformDeploymentEntity(entity)
	}
	return responses
}
