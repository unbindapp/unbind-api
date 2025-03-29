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
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *LogsService) GetLogs(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.LogQueryInput, send sse.Sender) error {
	team, project, environment, service, err := self.validatePermissionsAndParseInputs(ctx, requesterUserID, input)
	if err != nil {
		return err
	}

	// Build labels to select
	labels := map[string]string{}
	switch input.Type {
	case models.LogTypeTeam:
		labels["unbind-team"] = team.Name
	case models.LogTypeProject:
		labels["unbind-project"] = project.Name
	case models.LogTypeEnvironment:
		labels["unbind-environment"] = environment.Name
	case models.LogTypeService:
		labels["unbind-service"] = service.Name
	}

	// Build log options
	logOptions := k8s.LogOptions{}
	logOptions.Namespace = team.Namespace
	logOptions.Labels = labels
	logOptions.Previous = input.Previous
	logOptions.Tail = input.Tail
	logOptions.Timestamps = input.Timestamps
	logOptions.SearchPattern = input.SearchPattern
	logOptions.Follow = true

	// Parse 'since' duration
	var since time.Duration
	if input.Since != "" {
		since, err = time.ParseDuration(input.Since)
		if err != nil {
			return fmt.Errorf("invalid since duration: %w", err)
		}
	}
	logOptions.Since = since

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return err
	}

	// Create a channel for log events
	eventChan := make(chan k8s.LogEvent, 100)

	// Create a context with cancellation
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	pods, err := self.k8s.GetPodsByLabels(ctx, team.Namespace, labels, client)
	if err != nil {
		return fmt.Errorf("error getting pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found matching the specified criteria")
	}

	// Start streaming logs from each pod
	for _, pod := range pods.Items {
		serviceName, _ := pod.Labels["unbind-service"]
		service, err := self.repo.Service().GetByName(ctx, serviceName)
		if err != nil {
			return fmt.Errorf("error getting service: %w", err)
		}
		podName := pod.Name
		go func(podName string) {
			err := self.k8s.StreamPodLogs(streamCtx, podName, team.Namespace, logOptions, k8s.LogMetadata{
				ServiceID:     service.ID,
				EnvironmentID: service.Edges.Environment.ID,
				ProjectID:     service.Edges.Environment.Edges.Project.ID,
				TeamID:        service.Edges.Environment.Edges.Project.Edges.Team.ID,
			}, client, eventChan)
			if err != nil {
				// Send error as a log event
				select {
				case eventChan <- k8s.LogEvent{
					PodName: podName,
					Message: fmt.Sprintf("Error streaming logs: %v", err),
				}:
				case <-streamCtx.Done():
					return
				}
			}
		}(podName)
	}

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
