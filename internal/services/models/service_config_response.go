package models

import (
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// ServiceConfigResponse defines the configuration response for a service
type ServiceConfigResponse struct {
	GitBranch       *string               `json:"git_branch,omitempty"`
	Type            schema.ServiceType    `json:"type"`
	Builder         schema.ServiceBuilder `json:"builder"`
	Icon            string                `json:"icon"`
	Host            []v1.HostSpec         `json:"hosts,omitempty" nullable:"false"`
	Port            []schema.PortSpec     `json:"ports,omitempty" nullable:"false"`
	Replicas        int32                 `json:"replicas"`
	AutoDeploy      bool                  `json:"auto_deploy"`
	RunCommand      *string               `json:"run_command,omitempty"`
	Public          bool                  `json:"public"`
	Image           string                `json:"image,omitempty"`
	DatabaseVersion *string               `json:"database_version,omitempty"`
	DatabaseType    *string               `json:"database_type,omitempty"`
}

// TransformServiceConfigEntity transforms an ent.ServiceConfig entity into a ServiceConfigResponse
func TransformServiceConfigEntity(entity *ent.ServiceConfig) *ServiceConfigResponse {
	response := &ServiceConfigResponse{}
	if entity != nil {
		response = &ServiceConfigResponse{
			GitBranch:       entity.GitBranch,
			Type:            entity.Type,
			Builder:         entity.Builder,
			Icon:            entity.Icon,
			Host:            entity.Hosts,
			Port:            entity.Ports,
			Replicas:        entity.Replicas,
			AutoDeploy:      entity.AutoDeploy,
			RunCommand:      entity.RunCommand,
			Public:          entity.Public,
			Image:           entity.Image,
			DatabaseVersion: entity.DatabaseVersion,
			DatabaseType:    entity.Database,
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
