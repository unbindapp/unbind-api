package templates_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/pkg/templates"
)

type ListTemplatesInput struct {
	server.BaseAuthInput
	CategoryFilter templates.TemplateCategoryName `query:"category" required:"false" description:"Filter by category"`
}

type ListTemplatesResponse struct {
	Body struct {
		Data []*templates.TemplateList `json:"data" nullable:"false"`
	}
}

// ListTemplates handles GET /templates/list
func (self *HandlerGroup) ListTemplates(ctx context.Context, input *ListTemplatesInput) (*ListTemplatesResponse, error) {

	templateList, err := self.srv.TemplateProvider.ListTemplates(ctx, self.srv.Cfg.TemplateVersion, input.CategoryFilter)
	if err != nil {
		log.Errorf("failed to list templates: %v", err)
		return nil, huma.Error500InternalServerError("An unknown error occured")
	}

	return &ListTemplatesResponse{
		Body: struct {
			Data []*templates.TemplateList `json:"data" nullable:"false"`
		}{
			Data: []*templates.TemplateList{templateList},
		},
	}, nil
}
