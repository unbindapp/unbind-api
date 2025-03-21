package models

import "github.com/unbindapp/unbind-api/ent"

// ServiceConfigResponse defines the configuration response for a service
type ServiceConfigResponse struct {
	GitBranch  *string `json:"git_branch,omitempty"`
	Host       *string `json:"host,omitempty"`
	Port       *int    `json:"port,omitempty"`
	Replicas   int32   `json:"replicas"`
	AutoDeploy bool    `json:"auto_deploy"`
	RunCommand *string `json:"run_command,omitempty"`
	Public     bool    `json:"public"`
	Image      string  `json:"image,omitempty"`
}

// TransformServiceConfigEntity transforms an ent.ServiceConfig entity into a ServiceConfigResponse
func TransformServiceConfigEntity(entity *ent.ServiceConfig) *ServiceConfigResponse {
	response := &ServiceConfigResponse{}
	if entity != nil {
		response = &ServiceConfigResponse{
			GitBranch:  entity.GitBranch,
			Host:       entity.Host,
			Port:       entity.Port,
			Replicas:   entity.Replicas,
			AutoDeploy: entity.AutoDeploy,
			RunCommand: entity.RunCommand,
			Public:     entity.Public,
			Image:      entity.Image,
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
