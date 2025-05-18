package template_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type TemplateDeployInput struct {
	server.BaseAuthInput
	Body *models.TemplateDeployInput
}

type TemplateDeployResponse struct {
	Body struct {
		Data []*models.ServiceResponse `json:"data" required:"true"`
	}
}

func (self *HandlerGroup) DeployTemplate(ctx context.Context, input *TemplateDeployInput) (*TemplateDeployResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	// Deploy template
	services, err := self.srv.TemplateService.DeployTemplate(ctx, user.ID, bearerToken, input.Body)
	if err != nil {
		return nil, self.handleErr(err)
	}

	// Return response
	resp := &TemplateDeployResponse{}
	resp.Body.Data = services
	return resp, nil
}
