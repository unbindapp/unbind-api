package main

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/buildkitd.go"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

type Bootstrapper struct {
	repos                   repositories.RepositoriesInterface
	buildkitSettingsManager *buildkitd.BuildkitSettingsManager
}

func (self *Bootstrapper) Sync(ctx context.Context) error {
	var multierr error

	// Sync system settings
	multierror.Append(multierr, self.syncBuildkitdSettings(ctx))

	return multierr
}

// * Buildkit
func (self *Bootstrapper) syncBuildkitdSettings(ctx context.Context) error {
	// Get from kubernetes
	_, err := self.buildkitSettingsManager.GetOrCreateBuildkitConfig(ctx)
	if err != nil {
		log.Errorf("Failed to get buildkitd settings from kubernetes: %v", err)
		return err
	}

	// Get buildkit settings from DB
	settings, err := self.repos.System().GetBuildkitSettings(ctx)
	if err != nil && !ent.IsNotFound(err) {
		log.Errorf("Failed to get buildkit settings from DB: %v", err)
		return err
	}

	if ent.IsNotFound(err) {
		// Create record
		parallelism, err := self.buildkitSettingsManager.GetCurrentMaxParallelism(ctx)
		if err != nil {
			log.Errorf("Failed to get current max parallelism: %v", err)
			return err
		}
		// Get replicas
		replicas, err := self.buildkitSettingsManager.GetCurrentReplicas(ctx)
		if err != nil {
			log.Errorf("Failed to get current replicas: %v", err)
			return err
		}
		settings, err = self.repos.System().CreateBuildkitSettings(ctx, replicas, parallelism)
		if err != nil {
			log.Errorf("Failed to create buildkit settings in DB: %v", err)
			return err
		}

		log.Info("Created buildkitd settings", "replicas", replicas, "parallelism", parallelism)
		return nil
	}

	// Sync with kubernetes
	parallelism, err := self.buildkitSettingsManager.GetCurrentMaxParallelism(ctx)
	if err != nil {
		log.Errorf("Failed to get current max parallelism: %v", err)
		return err
	}

	// Get replicas
	replicas, err := self.buildkitSettingsManager.GetCurrentReplicas(ctx)
	if err != nil {
		log.Errorf("Failed to get current replicas: %v", err)
		return err
	}

	if settings.MaxParallelism != parallelism || settings.Replicas != replicas {
		// Update record
		settings, err = self.repos.System().UpdateBuildkitSettings(ctx, settings.ID, replicas, parallelism)
		if err != nil {
			log.Errorf("Failed to update buildkit settings in DB: %v", err)
			return err
		}

		log.Info("Updated buildkitd settings", "replicas", replicas, "parallelism", parallelism)
	}

	return nil
}
