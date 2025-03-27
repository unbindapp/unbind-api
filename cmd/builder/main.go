package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/database"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	"github.com/unbindapp/unbind-api/pkg/builder/builders"
	"github.com/unbindapp/unbind-api/pkg/builder/config"
	"github.com/unbindapp/unbind-api/pkg/builder/k8s"
	"gopkg.in/yaml.v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Warnf("Failed to load .env file: %v", err)
	}

	cfg := config.NewConfig()
	os.Setenv("BUILDKIT_HOST", cfg.BuildkitHost)

	// Setup database
	// Load database
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

	builder := builders.NewBuilder(cfg)
	k8s := k8s.NewK8SClient(cfg, cfg)

	// Parse secrets from env
	serializableSecrets := make(map[string]string)
	buildSecrets := make(map[string]string)
	if cfg.ServiceBuildSecrets != "" {
		if err := json.Unmarshal([]byte(cfg.ServiceBuildSecrets), &serializableSecrets); err != nil {
			log.Fatalf("Failed to parse secrets: %v", err)
		}

		// Convert back to map[string][]byte
		for k, v := range serializableSecrets {
			data, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				fmt.Printf("Error decoding secret %s: %v\n", k, err)
				continue
			}
			buildSecrets[k] = string(data)
		}
	}

	// Build with context
	var dockerImg, repoName string
	switch cfg.ServiceBuilder {
	case "railpack":
		dockerImg, repoName, err = builder.BuildWithRailpack(ctx, buildSecrets)
		if err != nil {
			log.Fatalf("Failed to build with railpack: %v", err)
		}
	default:
		log.Fatalf("Unknown builder: %s", cfg.ServiceBuilder)
	}

	crdName := cfg.ServiceName
	if crdName == "" {
		log.Warn("Service name not provided, using repository name")
		crdName = repoName
	}

	if dockerImg == "" {
		log.Error("Failed to build image to deploy!")
		os.Exit(1)
	}

	// Deploy to kubernetes with context
	createdCRD, serviceSpec, err := k8s.DeployImage(ctx, crdName, dockerImg)
	if err != nil {
		log.Fatalf("Failed to deploy image: %v", err)
	}

	// Update deployment metadata in the DB
	if _, err = repo.Deployment().AttachDeploymentMetadata(
		ctx,
		nil,
		cfg.ServiceDeploymentID,
		dockerImg,
		serviceSpec,
	); err != nil {
		log.Error("Failed to attach deployment metadata", "deployment_id", cfg.ServiceDeploymentID, "err", err)
	}

	// Update active deployment
	if err = repo.Service().SetCurrentDeployment(ctx, nil, cfg.ServiceID, cfg.ServiceDeploymentID); err != nil {
		log.Error("Failed to set current deployment", "service_id", cfg.ServiceID, "deployment_id", cfg.ServiceDeploymentID, "err", err)
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
