package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/infrastructure/database"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	webhooks_service "github.com/unbindapp/unbind-api/internal/services/webooks"
	"github.com/unbindapp/unbind-api/pkg/builder/builders"
	"github.com/unbindapp/unbind-api/pkg/builder/config"
	"github.com/unbindapp/unbind-api/pkg/builder/k8s"
	_ "go.uber.org/automaxprocs"
	"gopkg.in/yaml.v2"
)

func markDeploymentSuccessful(ctx context.Context, cfg *config.Config, webhooksService *webhooks_service.WebhooksService, tx repository.TxInterface, repo *repositories.Repositories, deploymentID uuid.UUID) error {
	_, err := repo.Deployment().MarkSucceeded(ctx, tx, deploymentID, time.Now())

	// Trigger webhook
	go func() {
		event := schema.WebhookProjectEventDeploymentSucceeded
		level := webhooks_service.WebhookLevelInfo

		// Get service with edges
		serviceID, _ := uuid.Parse(cfg.ServiceRef)
		service, err := repo.Service().GetByID(context.Background(), serviceID)
		if err != nil {
			log.Errorf("Failed to get service %s: %v", service.ID.String(), err)
			return
		}

		// Construct URL
		url, _ := utils.JoinURLPaths(cfg.ExternalUIUrl, service.Edges.Environment.Edges.Project.Edges.Team.ID.String(), "project", service.Edges.Environment.Edges.Project.ID.String(), "?environment="+service.EnvironmentID.String(), "&service="+service.ID.String(), "&deployment="+cfg.ServiceDeploymentID.String())
		data := webhooks_service.WebookData{
			Title:       "Deployment Succeeded",
			Url:         url,
			Description: fmt.Sprintf("A deployment has succeeded for %s", service.DisplayName),
			Fields: []webhooks_service.WebhookDataField{
				{
					Name:  "Service Type",
					Value: string(service.Edges.ServiceConfig.Type),
				},
				{
					Name:  "Environment",
					Value: service.Edges.Environment.Name,
				},
				{
					Name:  "Builder",
					Value: string(service.Edges.ServiceConfig.Builder),
				},
			},
		}

		if err := webhooksService.TriggerWebhooks(context.Background(), level, event, data); err != nil {
			log.Errorf("Failed to trigger webhook %s: %v", event, err)
		}
	}()

	return err
}

func markDeploymentFailed(ctx context.Context, cfg *config.Config, webhooksService *webhooks_service.WebhooksService, repo *repositories.Repositories, reason string, deploymentID uuid.UUID) error {
	_, err := repo.Deployment().MarkFailed(ctx, nil, deploymentID, reason, time.Now())
	// Trigger webhook
	go func() {
		event := schema.WebhookProjectEventDeploymentFailed
		level := webhooks_service.WebhookLevelError

		// Get service with edges
		serviceID, _ := uuid.Parse(cfg.ServiceRef)
		service, err := repo.Service().GetByID(context.Background(), serviceID)
		if err != nil {
			log.Errorf("Failed to get service %s: %v", service.ID.String(), err)
			return
		}

		// Construct URL
		url, _ := utils.JoinURLPaths(cfg.ExternalUIUrl, service.Edges.Environment.Edges.Project.Edges.Team.ID.String(), "project", service.Edges.Environment.Edges.Project.ID.String(), "?environment="+service.EnvironmentID.String(), "&service="+service.ID.String(), "&deployment="+cfg.ServiceDeploymentID.String())
		data := webhooks_service.WebookData{
			Title:       "Deployment Failed",
			Url:         url,
			Description: fmt.Sprintf("A build has failed for %s", service.DisplayName),
			Fields: []webhooks_service.WebhookDataField{
				{
					Name:  "Service Type",
					Value: string(service.Edges.ServiceConfig.Type),
				},
				{
					Name:  "Environment",
					Value: service.Edges.Environment.Name,
				},
				{
					Name:  "Builder",
					Value: string(service.Edges.ServiceConfig.Builder),
				},
				{
					Name:  "Error message",
					Value: reason,
				},
			},
		}

		if err := webhooksService.TriggerWebhooks(context.Background(), level, event, data); err != nil {
			log.Errorf("Failed to trigger webhook %s: %v", event, err)
		}
	}()
	return err
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	godotenv.Load()

	cfg := config.NewConfig()
	os.Setenv("BUILDKIT_HOST", cfg.BuildkitHost)

	log.Infof("Starting build...")
	log.Info("--------------")
	log.Infof("Input Parameters:")
	log.Infof(" - Service name: %s", cfg.ServiceName)
	if cfg.ServiceDatabaseType != "" {
		log.Infof(" - Using database type: %s", cfg.ServiceDatabaseType)
		log.Infof(" - Resource definition version: %s", cfg.ServiceDatabaseDefinitionVersion)
	} else if cfg.ServiceImage != "" {
		log.Infof(" - Using docker image: %s", cfg.ServiceImage)
	} else {
		log.Infof(" - Builder Type: %s", cfg.ServiceBuilder)
		if cfg.ServiceBuilder == schema.ServiceBuilderDocker {
			dockerfileDisplay := "Dockerfile"
			if cfg.ServiceDockerfilePath != "" {
				dockerfileDisplay = cfg.ServiceDockerfilePath
			}
			log.Infof(" - Dockerfile Path: %s", dockerfileDisplay)
			ctxDisplay := "."
			if cfg.ServiceDockerfileContext != "" {
				ctxDisplay = cfg.ServiceDockerfileContext
			}
			log.Infof(" - Dockerfile Context: %s", ctxDisplay)
		}
	}
	fmt.Printf("\n")

	serviceId, err := uuid.Parse(cfg.ServiceRef)
	if err != nil {
		log.Fatalf("Failed to parse service ID, ref must be a valid uuidv4: %v", err)
	}

	// Setup database
	dbConnInfo, err := database.GetSqlDbConn(cfg, false)
	if err != nil {
		log.Fatalf("Failed to get database connection info: %v", err)
	}
	// Initialize ent client
	db, err := database.NewEntClient(dbConnInfo)
	if err != nil {
		log.Fatalf("Failed to create ent client: %v", err)
	}
	repo := repositories.NewRepositories(db)
	webhooksService := webhooks_service.NewWebhooksService(repo)

	// Trigger webhook
	go func() {
		event := schema.WebhookProjectEventDeploymentBuilding
		level := webhooks_service.WebhookLevelInfo

		// Get service with edges
		serviceID, _ := uuid.Parse(cfg.ServiceRef)
		service, err := repo.Service().GetByID(context.Background(), serviceID)
		if err != nil {
			log.Errorf("Failed to get service %s: %v", service.ID.String(), err)
			return
		}

		// Construct URL
		url, _ := utils.JoinURLPaths(cfg.ExternalUIUrl, service.Edges.Environment.Edges.Project.Edges.Team.ID.String(), "project", service.Edges.Environment.Edges.Project.ID.String(), "?environment="+service.EnvironmentID.String(), "&service="+service.ID.String(), "&deployment="+cfg.ServiceDeploymentID.String())
		data := webhooks_service.WebookData{
			Title:       "Deployment Building",
			Url:         url,
			Description: fmt.Sprintf("A build has started for %s", service.DisplayName),
			Fields: []webhooks_service.WebhookDataField{
				{
					Name:  "Service Type",
					Value: string(service.Edges.ServiceConfig.Type),
				},
				{
					Name:  "Environment",
					Value: service.Edges.Environment.Name,
				},
				{
					Name:  "Builder",
					Value: string(service.Edges.ServiceConfig.Builder),
				},
			},
		}

		if err := webhooksService.TriggerWebhooks(context.Background(), level, event, data); err != nil {
			log.Errorf("Failed to trigger webhook %s: %v", event, err)
		}
	}()

	builder := builders.NewBuilder(cfg)
	k8s := k8s.NewK8SClient(cfg, cfg)

	var dockerImg string

	// We can bypass any build step if the image is already provided
	if cfg.ServiceImage != "" || cfg.ServiceType == schema.ServiceTypeDatabase {
		dockerImg = cfg.ServiceImage
	} else {
		// Parse secrets from env
		serializableSecrets := make(map[string]string)
		buildSecrets := make(map[string]string)
		if cfg.ServiceBuildSecrets != "" {
			if err := json.Unmarshal([]byte(cfg.ServiceBuildSecrets), &serializableSecrets); err != nil {
				if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("failed to unmarshal secrets %v", err), cfg.ServiceDeploymentID); err != nil {
					log.Errorf("Failed to mark deployment as failed: %v", err)
				}
				log.Fatalf("Failed to parse secrets: %v", err)
			}

			// Convert back to map[string][]byte
			for k, v := range serializableSecrets {
				data, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					log.Warnf("Error decoding secret %s: %v\n", k, err)
					continue
				}
				buildSecrets[k] = string(data)
			}
		}

		// Build with context
		switch cfg.ServiceBuilder {
		case schema.ServiceBuilderRailpack:
			dockerImg, _, err = builder.BuildWithRailpack(ctx, buildSecrets)
			if err != nil {
				if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("failed railpack build %v", err), cfg.ServiceDeploymentID); err != nil {
					log.Errorf("Failed to mark deployment as failed: %v", err)
				}
				log.Fatalf("Failed to build with railpack: %v", err)
			}
		case schema.ServiceBuilderDocker:
			dockerImg, _, err = builder.BuildDockerfile(ctx, buildSecrets)
			if err != nil {
				if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("failed docker build %v", err), cfg.ServiceDeploymentID); err != nil {
					log.Errorf("Failed to mark deployment as failed: %v", err)
				}
				log.Fatalf("Failed to build with docker: %v", err)
			}
		default:
			if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("received request with unknown builder: %s", cfg.ServiceBuilder), cfg.ServiceDeploymentID); err != nil {
				log.Errorf("Failed to mark deployment as failed: %v", err)
			}
			log.Fatalf("Unknown builder: %s", cfg.ServiceBuilder)
		}
	}

	crdName := cfg.ServiceName
	if crdName == "" {
		log.Fatal("Service name not provided, cannot deploy")
	}

	// Database doesn't need a build, so bypass
	if dockerImg == "" && cfg.ServiceType != schema.ServiceTypeDatabase {
		if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("no output image generated"), cfg.ServiceDeploymentID); err != nil {
			log.Errorf("Failed to mark deployment as failed: %v", err)
		}
		log.Error("Failed to build image to deploy!")
		os.Exit(1)
	}

	// Deploy to kubernetes with context
	createdCRD, serviceSpec, err := k8s.DeployImage(ctx, crdName, dockerImg)
	if err != nil {
		if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("failed to deploy image %v", err), cfg.ServiceDeploymentID); err != nil {
			log.Errorf("Failed to mark deployment as failed: %v", err)
		}
		log.Fatalf("Failed to deploy image: %v", err)
	}

	// Update deployment metadata in the DB
	if err := repo.WithTx(ctx, func(tx repository.TxInterface) error {
		if _, err = repo.Deployment().AttachDeploymentMetadata(
			ctx,
			tx,
			cfg.ServiceDeploymentID,
			dockerImg,
			serviceSpec,
		); err != nil {
			log.Error("Failed to attach deployment metadata", "deployment_id", cfg.ServiceDeploymentID, "err", err)
		}

		// Update active deployment
		if err = repo.Service().SetCurrentDeployment(ctx, tx, serviceId, cfg.ServiceDeploymentID); err != nil {
			log.Error("Failed to set current deployment", "service_id", serviceId, "deployment_id", cfg.ServiceDeploymentID, "err", err)
		}

		if err = markDeploymentSuccessful(ctx, cfg, webhooksService, tx, repo, cfg.ServiceDeploymentID); err != nil {
			log.Error("Failed to mark deployment as successful", "deployment_id", cfg.ServiceDeploymentID, "err", err)
		}
		return nil
	}); err != nil {
		if err := markDeploymentFailed(ctx, cfg, webhooksService, repo, fmt.Sprintf("failed to update deployment metadata %v", err), cfg.ServiceDeploymentID); err != nil {
			log.Errorf("Failed to mark deployment as failed: %v", err)
		}
		log.Fatalf("Failed to update deployment metadata: %v", err)
	}

	// Pretty print the CRD as YAML
	crdYAML, err := yaml.Marshal(createdCRD)
	if err != nil {
		log.Errorf("Failed to marshal CRD to YAML: %v", err)
		log.Infof("Created CRD: %v", createdCRD) // Fallback to default printing
	} else {
		log.Infof("Created CRD:\n%s", string(crdYAML))
	}
}
