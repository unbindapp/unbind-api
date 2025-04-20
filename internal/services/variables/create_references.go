package variables_service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *VariablesService) CreateVariableReference(ctx context.Context, requesterUserID uuid.UUID, input *models.CreateVariableReferenceInput) (*models.VariableReferenceResponse, error) {
	// ! TODO - we're going to need to change all of our permission checks to filter not reject
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.TargetServiceID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Validate the template with sources
	if err := ValidateCreateVariableReferenceInput(input); err != nil {
		return nil, err
	}

	// Ensure service exists
	_, err := self.repo.Service().GetByID(ctx, input.TargetServiceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	resp, err := self.repo.Variables().CreateReference(ctx, input)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Variable reference already exists")
		}
		return nil, err
	}

	return models.TransformVariableReferenceResponseEntity(resp), nil
}

func ValidateCreateVariableReferenceInput(input *models.CreateVariableReferenceInput) error {
	if len(input.Sources) == 0 {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "At least one source is required")
	}

	template := input.ValueTemplate

	// Track which sources have been referenced
	sourcesReferenced := make(map[string]bool)
	for _, source := range input.Sources {
		sourcesReferenced[source.Name+"."+source.Key] = false
	}

	// Find all occurrences of ${...} in the template that aren't escaped
	// We'll use regexp to find all instances, then check if they're escaped
	re := regexp.MustCompile(`\${([^}]+)}`)
	matches := re.FindAllStringSubmatchIndex(template, -1)

	for _, match := range matches {
		// Check if this is an escaped instance (preceded by \)
		if match[0] > 0 && template[match[0]-1] == '\\' {
			continue // Skip escaped instances
		}

		// Extract the source reference (name.key)
		reference := template[match[2]:match[3]]

		// Mark this source as referenced if it exists in our sources
		if _, exists := sourcesReferenced[reference]; exists {
			sourcesReferenced[reference] = true
		}
	}

	// Check if all sources have been referenced
	var missingReferences []string
	for sourceRef, referenced := range sourcesReferenced {
		if !referenced {
			missingReferences = append(missingReferences, sourceRef)
		}
	}

	if len(missingReferences) > 0 {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("template is missing references to the following sources: %s",
			strings.Join(missingReferences, ", ")))
	}

	return nil
}
