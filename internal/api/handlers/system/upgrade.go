package system_handler

import (
	"context"
	"slices"
	"sort"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"golang.org/x/mod/semver"
)

func (self *HandlerGroup) CheckPermissions(ctx context.Context, requesterUserID uuid.UUID) error {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Team editor can create projects
		{
			Action:       schema.ActionAdmin,
			ResourceType: schema.ResourceTypeSystem,
		},
	}

	if err := self.srv.Repository.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return huma.Error403Forbidden("You are not authorized to perform this action")
	}

	return nil
}

// * Check for updates
type UpgradeCheckResponse struct {
	Body struct {
		HasUpgradeAvailable bool     `json:"has_upgrade_available"`
		AvailableVersions   []string `json:"available_versions"`
		CurrentVersion      string   `json:"current_version"`
	}
}

func (self *HandlerGroup) CheckForUpdates(ctx context.Context, input *server.BaseAuthInput) (*UpgradeCheckResponse, error) {
	// Get requester
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	// Check permissions
	if err := self.CheckPermissions(ctx, user.ID); err != nil {
		return nil, err
	}

	// Get all available versions
	allUpdates, err := self.srv.UpgradeManager.CheckForUpdates(ctx)
	if err != nil {
		// Log the error but return empty updates instead of error
		log.Errorf("Failed to check for updates: %v", err)
		resp := &UpgradeCheckResponse{}
		resp.Body.HasUpgradeAvailable = false
		resp.Body.AvailableVersions = []string{}
		resp.Body.CurrentVersion = self.srv.UpgradeManager.CurrentVersion
		return resp, nil
	}

	// Filter to only show versions that can be upgraded to in sequence
	availableUpdates := make([]string, 0)
	currentVersion := self.srv.UpgradeManager.CurrentVersion

	// Sort versions to ensure we process them in order
	sort.Slice(allUpdates, func(i, j int) bool {
		return semver.Compare(allUpdates[i], allUpdates[j]) < 0
	})

	// Keep track of the current version we're checking from
	checkVersion := currentVersion

	// Build the upgrade path
	for {
		nextVersion, err := self.srv.UpgradeManager.GetNextAvailableVersion(ctx, checkVersion)
		if err != nil || nextVersion == "" {
			// No more versions available
			break
		}

		// Add the next version to our list
		availableUpdates = append(availableUpdates, nextVersion)

		// Update the version we're checking from
		checkVersion = nextVersion
	}

	resp := &UpgradeCheckResponse{}
	resp.Body.HasUpgradeAvailable = len(availableUpdates) > 0
	resp.Body.AvailableVersions = availableUpdates
	resp.Body.CurrentVersion = currentVersion

	return resp, nil
}

// * Apply update
type UpgradeApplyInput struct {
	server.BaseAuthInput
	Body struct {
		TargetVersion string `json:"target_version"`
	}
}

type UpgradeApplyResponse struct {
	Body struct {
		Started bool `json:"started"`
	}
}

func (self *HandlerGroup) ApplyUpdate(ctx context.Context, input *UpgradeApplyInput) (*UpgradeApplyResponse, error) {
	// Get requester
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	// Check permissions
	if err := self.CheckPermissions(ctx, user.ID); err != nil {
		return nil, err
	}

	// Validate version is available
	availableUpdates, err := self.srv.UpgradeManager.CheckForUpdates(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check for updates: " + err.Error())
	}

	if !slices.Contains(availableUpdates, input.Body.TargetVersion) {
		return nil, huma.Error400BadRequest("Target version is not available")
	}

	// ! Temporarily disabling upgrade
	// Apply update
	// err = self.srv.UpgradeManager.UpgradeToVersion(ctx, input.Body.TargetVersion)
	// if err != nil {
	// 	return nil, huma.Error500InternalServerError("Failed to apply update: " + err.Error())
	// }

	resp := &UpgradeApplyResponse{}
	resp.Body.Started = true

	return resp, nil
}

// * Get upgrade status
type UpgradeStatusInput struct {
	ExpectedVersion string `query:"expected_version"`
}

type UpgradeStatusResponse struct {
	Body struct {
		Ready bool `json:"ready"`
	}
}

func (self *HandlerGroup) GetUpgradeStatus(ctx context.Context, input *UpgradeStatusInput) (*UpgradeStatusResponse, error) {
	// Get requester
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	// Check permissions
	if err := self.CheckPermissions(ctx, user.ID); err != nil {
		return nil, err
	}

	// Get upgrade status
	ready, err := self.srv.UpgradeManager.CheckDeploymentsReady(ctx, input.ExpectedVersion)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get upgrade status: " + err.Error())
	}

	resp := &UpgradeStatusResponse{}
	resp.Body.Ready = ready

	return resp, nil
}
