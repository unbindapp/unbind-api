package logs_service

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
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
		podName := pod.Name
		go func(podName string) {
			err := self.k8s.StreamPodLogs(streamCtx, podName, team.Namespace, logOptions, eventChan)
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
		// Has permission to read system resources
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to read this specific team
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
	}

	// Depends on type of logs for other permissions
	if input.Type == models.LogTypeProject ||
		input.Type == models.LogTypeEnvironment ||
		input.Type == models.LogTypeService {
		if input.ProjectID == uuid.Nil {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project ID is required")
		}
		// Has permission to read projects
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		})
		// Has permission to read this specific project
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   input.ProjectID.String(),
		})
	}

	if input.Type == models.LogTypeEnvironment ||
		input.Type == models.LogTypeService {
		if input.EnvironmentID == uuid.Nil {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment ID is required")
		}
		// Has permission to read environments
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   "*",
		})
		// Has permission to read this specific environment
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID.String(),
		})
	}

	if input.Type == models.LogTypeService {
		if input.ServiceID == uuid.Nil {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Service ID is required")
		}
		// Has permission to read services
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   "*",
		})
		// Has permission to read this specific service
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   input.ServiceID.String(),
		})
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
