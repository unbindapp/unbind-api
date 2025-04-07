package templates_handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/pkg/templates"
)

type GetTemplateMetadataInput struct {
	server.BaseAuthInput
	Category templates.TemplateCategoryName `query:"category" required:"true" description:"Category of the template"`
	Name     string                         `query:"name" required:"true" description:"Name of the template"`
	Version  string                         `query:"version" required:"false" description:"Version of the templates release"`
}

type GetTemplateResponse struct {
	Body struct {
		Data *templates.Template `json:"data" nullable:"false"`
	}
}

// ListTemplates handles GET /templates/get
func (self *HandlerGroup) GetTemplateMetadata(ctx context.Context, input *GetTemplateMetadataInput) (*GetTemplateResponse, error) {
	version := self.srv.Cfg.TemplateVersion
	if input.Version != "" {
		version = input.Version
	}

	template, err := self.srv.TemplateProvider.FetchTemplate(
		ctx,
		version,
		input.Category,
		input.Name,
	)
	if err != nil {
		if errors.Is(err, templates.ErrTemplateNotFound) {
			return nil, huma.Error404NotFound("Template not found")
		}
		log.Errorf("failed to get templates: %v", err)
		return nil, huma.Error500InternalServerError("An unknown error occured")
	}

	response := &GetTemplateResponse{}
	response.Body.Data = template
	return response, nil
}
