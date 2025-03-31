package logs_service

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *LogsService) GetLogs(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.LogQueryInput, send sse.Sender) error {
	team, project, environment, service, err := self.validatePermissionsAndParseInputs(ctx, requesterUserID, input)
	if err != nil {
		return err
	}

	// Build labels to select
	var label loki.LokiLabelName
	var labelValue string
	switch input.Type {
	case models.LogTypeTeam:
		label = loki.LokiLabelTeam
		labelValue = team.ID.String()
	case models.LogTypeProject:
		label = loki.LokiLabelProject
		labelValue = project.ID.String()
	case models.LogTypeEnvironment:
		label = loki.LokiLabelEnvironment
		labelValue = environment.ID.String()
	case models.LogTypeService:
		label = loki.LokiLabelService
		labelValue = service.ID.String()
	}

	// Parse 'since' duration
	var since time.Duration
	if input.Since != "" {
		since, err = time.ParseDuration(input.Since)
		if err != nil {
			return fmt.Errorf("invalid since duration: %w", err)
		}
	}

	// Create a channel for log events
	eventChan := make(chan loki.LogEvents, 100)

	// Create a context with cancellation
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create loki options
	lokiLogOptions := loki.LokiLogOptions{
		Label:      label,
		LabelValue: labelValue,
		Limit:      int(input.Tail),
		RawFilter:  input.Filters,
		Since:      since,
	}

	// Start a single stream for all pods
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered from panic in log streaming goroutine: %v", r)
			}
		}()

		err := self.lokiQuerier.StreamLokiPodLogs(streamCtx, lokiLogOptions, eventChan)
		if err != nil {
			// Wrap the send in a select with context check to avoid sending on canceled contexts
			select {
			case <-streamCtx.Done():
				return
			default:
				send.Data(loki.LogsError{
					Code:    500,
					Message: fmt.Sprintf("Error streaming logs: %v", err),
				})
			}
		}
	}()

	// Send events to the client
	for {
		select {
		case event := <-eventChan:
			send.Data(event)
		case <-ctx.Done():
			return nil
		}
	}
}

func (self *LogsService) validatePermissionsAndParseInputs(ctx context.Context, requesterUserID uuid.UUID, input *models.LogQueryInput) (*ent.Team, *ent.Project, *ent.Environment, *ent.Service, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		//Can read team, project, environmnent, or service depending on inputs
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   input.TeamID,
		},
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   input.ProjectID,
		},
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID,
		},
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.ServiceID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, nil, nil, nil, err
	}

	// Get namespace
	team, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, nil, nil, nil, err
	}

	// Get project
	var project *ent.Project
	if input.Type == models.LogTypeProject ||
		input.Type == models.LogTypeEnvironment ||
		input.Type == models.LogTypeService {
		// validate project ID
		project, err = self.repo.Project().GetByID(ctx, input.ProjectID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
			}
			return nil, nil, nil, nil, err
		}
		if project.TeamID != input.TeamID {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found in this team")
		}
	}

	// Get environment
	var environment *ent.Environment
	if input.Type == models.LogTypeEnvironment ||
		input.Type == models.LogTypeService {
		// validate environment ID
		environment, err = self.repo.Environment().GetByID(ctx, input.EnvironmentID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
			}
			return nil, nil, nil, nil, err
		}
		if environment.ProjectID != input.ProjectID {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found in this project")
		}
	}

	// Get service
	var service *ent.Service
	if input.Type == models.LogTypeService {
		service, err = self.repo.Service().GetByID(ctx, input.ServiceID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
			}
			return nil, nil, nil, nil, err
		}
		if service.EnvironmentID != input.EnvironmentID {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found in this environment")
		}
	}

	return team, project, environment, service, nil
}
