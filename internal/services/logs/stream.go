package logs_service

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *LogsService) StreamLogs(ctx context.Context, requesterUserID uuid.UUID, input *models.LogStreamInput, send sse.Sender) error {
	team, project, environment, service, err := self.validatePermissionsAndParseInputs(ctx, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID)
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
	case models.LogTypeDeployment:
		// get deployment
		deployment, err := self.repo.Deployment().GetByID(ctx, input.DeploymentID)
		if err != nil {
			if ent.IsNotFound(err) {
				return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Deployment not found")
			}
			return err
		}
		label = loki.LokiLabelDeployment
		labelValue = deployment.ID.String()
	}

	// Parse 'since' duration
	var since time.Duration
	if input.Since != "" {
		since, err = time.ParseDuration(input.Since)
		if err != nil {
			return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "invalid since duration")
		}
	}

	// Create a channel for log events
	eventChan := make(chan loki.LogEvents, 100)

	// Create a context with cancellation
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create loki options
	lokiLogOptions := loki.LokiLogStreamOptions{
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
				send.Data(loki.LogEvents{
					MessageType:  loki.LogEventsMessageTypeError,
					ErrorMessage: fmt.Sprintf("Error streaming logs: %v", err),
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
