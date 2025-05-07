package templates_service

import (
	"context"

	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *TemplatesService) GetAvailable(ctx context.Context) ([]*models.TemplateResponse, error) {
	//  No special permission checks for reading these

	templates, err := self.repo.Template().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Transform the entities into response models
	return models.TransformTemplateEntities(templates), nil
}
