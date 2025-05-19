package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// ServiceConfigResponse defines the configuration response for a service
type ServiceConfigResponse struct {
	GitBranch      *string               `json:"git_branch,omitempty"`
	GitTag         *string               `json:"git_tag,omitempty"`
	Builder        schema.ServiceBuilder `json:"builder"`
	Icon           string                `json:"icon"`
	Host           []v1.HostSpec         `json:"hosts,omitempty" nullable:"false"`
	Port           []schema.PortSpec     `json:"ports,omitempty" nullable:"false"`
	Replicas       int32                 `json:"replicas"`
	AutoDeploy     bool                  `json:"auto_deploy"`
	InstallCommand *string               `json:"install_command,omitempty"`
	BuildCommand   *string               `json:"build_command,omitempty"`
	RunCommand     *string               `json:"run_command,omitempty"`
	IsPublic       bool                  `json:"is_public"`
	Image          string                `json:"image,omitempty"`
	// For backups
	S3BackupSourceID     *uuid.UUID `json:"s3_backup_source_id,omitempty"`
	S3BackupBucket       *string    `json:"s3_backup_bucket,omitempty"`
	BackupSchedule       string     `json:"backup_schedule"`
	BackupRetentionCount int        `json:"backup_retention_count"`
	// Volume
	Volumes []schema.ServiceVolume `json:"volumes" nullable:"false"`
	// Security context
	SecurityContext *schema.SecurityContext `json:"security_context,omitempty"`
	// Health check
	HealthCheck *schema.HealthCheck `json:"health_check,omitempty"`
	// Variable Volume Mounts
	VariableMounts []*schema.VariableMount `json:"variable_mounts" nullable:"false"`
	// Protected variables
	ProtectedVariables []string `json:"protected_variables" nullable:"false"`
	// Init containers
	InitContainers []*schema.InitContainer `json:"init_containers" nullable:"false"`
}

// TransformServiceConfigEntity transforms an ent.ServiceConfig entity into a ServiceConfigResponse
func TransformServiceConfigEntity(entity *ent.ServiceConfig) *ServiceConfigResponse {
	response := &ServiceConfigResponse{
		ProtectedVariables: []string{},
	}
	if entity != nil {
		if entity.VariableMounts == nil {
			entity.VariableMounts = []*schema.VariableMount{}
		}
		if entity.Volumes == nil {
			entity.Volumes = []schema.ServiceVolume{}
		}
		if entity.InitContainers == nil {
			entity.InitContainers = []*schema.InitContainer{}
		}
		response = &ServiceConfigResponse{
			GitBranch:            entity.GitBranch,
			GitTag:               entity.GitTag,
			Builder:              entity.Builder,
			Icon:                 entity.Icon,
			Host:                 entity.Hosts,
			Port:                 entity.Ports,
			Replicas:             entity.Replicas,
			AutoDeploy:           entity.AutoDeploy,
			InstallCommand:       entity.InstallCommand,
			BuildCommand:         entity.BuildCommand,
			RunCommand:           entity.RunCommand,
			IsPublic:             entity.IsPublic,
			Image:                entity.Image,
			S3BackupSourceID:     entity.S3BackupSourceID,
			S3BackupBucket:       entity.S3BackupBucket,
			BackupSchedule:       entity.BackupSchedule,
			BackupRetentionCount: entity.BackupRetentionCount,
			Volumes:              entity.Volumes,
			SecurityContext:      entity.SecurityContext,
			HealthCheck:          entity.HealthCheck,
			VariableMounts:       entity.VariableMounts,
			ProtectedVariables:   entity.ProtectedVariables,
			InitContainers:       entity.InitContainers,
		}
		if response.ProtectedVariables == nil {
			response.ProtectedVariables = []string{}
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
