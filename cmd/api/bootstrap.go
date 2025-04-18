package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/infrastructure/buildkitd"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

type Bootstrapper struct {
	cfg                     *config.Config
	repos                   repositories.RepositoriesInterface
	kubeClient              *k8s.KubeClient
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

// * Bootstrap the registry, if needed
func (self *Bootstrapper) bootstrapRegistry(ctx context.Context) error {
	regCount, err := self.repos.Ent().Registry.Query().Count(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return fmt.Errorf("failed to query registry count: %w", err)
		}
	}
	if regCount == 0 {
		// Create initial registry
		if self.cfg.BootstrapContainerRegistryHost == "" {
			return fmt.Errorf("BOOTSTRAP_CONTAINER_REGISTRY_HOST is empty")
		}

		if err := self.repos.WithTx(ctx, func(tx repository.TxInterface) error {
			// Create credentials first if needed
			if self.cfg.BootstrapContainerRegistryUser != "" {
				if self.cfg.BootstrapContainerRegistryPassword == "" {
					return fmt.Errorf("BOOTSTRAP_CONTAINER_REGISTRY_PASSWORD is empty")
				}

				secretName, err := utils.GenerateSlug(self.cfg.BootstrapContainerRegistryHost)
				if err != nil {
					return fmt.Errorf("failed to generate slug for registry secret: %w", err)
				}

				// Create registry
				_, err = self.repos.System().CreateRegistry(ctx, tx, self.cfg.BootstrapContainerRegistryHost, &secretName, true)
				if err != nil {
					return fmt.Errorf("failed to create registry: %w", err)
				}

				// Get teams so we can copy secret
				teams, err := self.repos.Ent().Team.Query().All(ctx)
				if err != nil {
					return fmt.Errorf("failed to query teams: %w", err)
				}

				// Create credentials second, so if it fails the above rolls back
				_, err = self.kubeClient.CreateMultiRegistryCredentials(
					ctx,
					secretName,
					self.cfg.SystemNamespace,
					[]k8s.RegistryCredential{
						{
							RegistryURL: self.cfg.BootstrapContainerRegistryHost,
							Username:    self.cfg.BootstrapContainerRegistryUser,
							Password:    self.cfg.BootstrapContainerRegistryPassword,
						},
					},
					self.kubeClient.GetInternalClient(),
				)
				if err != nil {
					return fmt.Errorf("failed to create registry credentials: %w", err)
				}

				for _, team := range teams {
					_, err = self.kubeClient.CopySecret(ctx, secretName, self.cfg.SystemNamespace, team.Namespace, self.kubeClient.GetInternalClient())
					if err != nil {
						return fmt.Errorf("failed to copy registry credentials to team namespace: %w", err)
					}
				}
			}

			log.Warn("Creating initial registry without credentials")

			_, err = self.repos.System().CreateRegistry(ctx, tx, self.cfg.BootstrapContainerRegistryHost, nil, true)
			if err != nil {
				return fmt.Errorf("failed to create registry: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("failed to bootstrap registry: %w", err)
		}
	}
	return nil
}
