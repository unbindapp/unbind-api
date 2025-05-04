package k8s

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"k8s.io/apimachinery/pkg/api/errors"
)

// SyncDatabaseSecrets syncs all database secrets with the operator logic
func (self *KubeClient) SyncDatabaseSecrets(ctx context.Context) error {
	// Get all database services
	databaseServices, err := self.repo.Service().GetDatabases(ctx)
	if err != nil {
		return err
	}

	for _, service := range databaseServices {
		if err := self.SyncDatabaseSecretForService(ctx, service); err != nil {
			log.Errorf("Failed to sync secret for service %s: %v", service.ID, err)
			// Continue with other services even if one fails
			continue
		}
	}
	return nil
}

// SyncDatabaseSecretForServiceID syncs the database secret for a specific service ID
func (self *KubeClient) SyncDatabaseSecretForServiceID(ctx context.Context, serviceID uuid.UUID) error {
	// Get the specific service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		return fmt.Errorf("failed to get service %s: %w", serviceID, err)
	}

	return self.SyncDatabaseSecretForService(ctx, service)
}

// SyncDatabaseSecretForService syncs the database secret for a specific service
func (self *KubeClient) SyncDatabaseSecretForService(ctx context.Context, service *ent.Service) error {
	if service.Type != schema.ServiceTypeDatabase {
		return nil
	}

	// Validate service has necessary relationships
	if service.Database == nil || service.Edges.Environment == nil ||
		service.Edges.Environment.Edges.Project == nil ||
		service.Edges.Environment.Edges.Project.Edges.Team == nil {
		return fmt.Errorf("service %s does not have a database or environment or project or team", service.ID)
	}

	namespace := service.Edges.Environment.Edges.Project.Edges.Team.Namespace

	// Get the existing secret
	secret, err := self.GetSecret(ctx, service.KubernetesSecret, namespace, self.GetInternalClient())
	if err != nil {
		return fmt.Errorf("failed to get secret %s in namespace %s: %w", service.KubernetesSecret, namespace, err)
	}

	username := string(secret.Data["DATABASE_USERNAME"])
	password := string(secret.Data["DATABASE_PASSWORD"])
	existingUrl := string(secret.Data["DATABASE_URL"])

	// For postgres, we can sync username and password if they are empty
	if *service.Database == "postgres" && (username == "" || password == "") {
		zalandoSecretName := fmt.Sprintf("postgres.%s.credentials.postgresql.acid.zalan.do", service.Name)
		zalandoSecret, err := self.GetSecret(ctx, zalandoSecretName, namespace, self.GetInternalClient())
		if err != nil {
			if errors.IsNotFound(err) {
				return fmt.Errorf("secret %s in namespace %s not found: %w", zalandoSecretName, namespace, err)
			}
			return fmt.Errorf("failed to get secret %s in namespace %s: %w", zalandoSecretName, namespace, err)
		}
		username = string(zalandoSecret.Data["username"])
		password = string(zalandoSecret.Data["password"])
	}

	if username == "" || password == "" {
		return fmt.Errorf("secret %s in namespace %s does not have username or password", service.KubernetesSecret, namespace)
	}

	// Set database URL for database type
	var url string
	switch *service.Database {
	case "postgres":
		url = fmt.Sprintf("postgresql://%s:%s@%s.%s:%d/postgres?sslmode=disable", username, password, service.KubernetesName, namespace, 5432)
	case "redis":
		url = fmt.Sprintf("redis://%s:%s@%s-headless.%s:%d", "default", password, service.KubernetesName, namespace, 6379)
	case "mysql":
		url = fmt.Sprintf("mysql://%s:%s@moco-%s.%s:%d/%s", username, password, service.KubernetesName, namespace, 3306, "moco")
	}

	if existingUrl != url {
		secrets := map[string][]byte{
			"DATABASE_USERNAME": []byte(username),
			"DATABASE_PASSWORD": []byte(password),
		}
		if url != "" {
			secrets["DATABASE_URL"] = []byte(url)
		}
		// Sync secret
		_, err = self.UpsertSecretValues(ctx, secret.Name, namespace, secrets, self.GetInternalClient())
		if err != nil {
			return fmt.Errorf("failed to update secret %s in namespace %s: %w", secret.Name, namespace, err)
		}
	}

	return nil
}
