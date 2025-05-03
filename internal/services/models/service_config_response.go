package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// ServiceConfigResponse defines the configuration response for a service
type ServiceConfigResponse struct {
	GitBranch  *string               `json:"git_branch,omitempty"`
	GitTag     *string               `json:"git_tag,omitempty"`
	Builder    schema.ServiceBuilder `json:"builder"`
	Icon       string                `json:"icon"`
	Host       []v1.HostSpec         `json:"hosts,omitempty" nullable:"false"`
	Port       []schema.PortSpec     `json:"ports,omitempty" nullable:"false"`
	Replicas   int32                 `json:"replicas"`
	AutoDeploy bool                  `json:"auto_deploy"`
	RunCommand *string               `json:"run_command,omitempty"`
	IsPublic   bool                  `json:"is_public"`
	Image      string                `json:"image,omitempty"`
	// For backups
	S3BackupSourceID *uuid.UUID `json:"s3_backup_source_id,omitempty"`
	S3BackupBucket   *string    `json:"s3_backup_bucket,omitempty"`
}

// TransformServiceConfigEntity transforms an ent.ServiceConfig entity into a ServiceConfigResponse
func TransformServiceConfigEntity(entity *ent.ServiceConfig) *ServiceConfigResponse {
	response := &ServiceConfigResponse{}
	if entity != nil {
		response = &ServiceConfigResponse{
			GitBranch:        entity.GitBranch,
			GitTag:           entity.GitTag,
			Builder:          entity.Builder,
			Icon:             entity.Icon,
			Host:             entity.Hosts,
			Port:             entity.Ports,
			Replicas:         entity.Replicas,
			AutoDeploy:       entity.AutoDeploy,
			RunCommand:       entity.RunCommand,
			IsPublic:         entity.IsPublic,
			Image:            entity.Image,
			S3BackupSourceID: entity.S3BackupSourceID,
			S3BackupBucket:   entity.S3BackupBucket,
		}
	}
	return response
}

// TransformServiceConfigEntities transforms a slice of ent.ServiceConfig entities into a slice of ServiceConfigResponse
func TransformServiceConfigEntities(entities []*ent.ServiceConfig) []*ServiceConfigResponse {
	responses := make([]*ServiceConfigResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformServiceConfigEntity(entity)
	}
	return responses
}
