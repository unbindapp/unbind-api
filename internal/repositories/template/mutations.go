package template_repo

import (
	"context"

	"github.com/hashicorp/go-multierror"
	entTemplate "github.com/unbindapp/unbind-api/ent/template"
	"github.com/unbindapp/unbind-api/pkg/templates"
)

func (self *TemplateRepository) UpsertPredefinedTemplates(ctx context.Context) (result error) {
	templates := templates.NewTemplater(nil).AvailableTemplates()

	for _, template := range templates {
		err := self.base.DB.Template.Create().
			SetName(template.Name).
			SetDisplayRank(template.DisplayRank).
			SetIcon(template.Icon).
			SetDescription(template.Description).
			SetKeywords(template.Keywords).
			SetResourceRecommendations(template.ResourceRecommendations).
			SetVersion(template.Version).
			SetDefinition(*template).
			SetImmutable(true).
			OnConflictColumns(entTemplate.FieldName, entTemplate.FieldVersion).
			DoNothing().
			Exec(ctx)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
