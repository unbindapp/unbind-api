package system_repo

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/registry"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *SystemRepository) CreateRegistry(ctx context.Context, tx repository.TxInterface, host string, kubernetesSecret *string, isDefault bool) (*ent.Registry, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	// Create registry
	return db.Registry.Create().
		SetHost(host).
		SetNillableKubernetesSecret(kubernetesSecret).
		SetIsDefault(isDefault).
		Save(ctx)
}

func (self *SystemRepository) GetDefaultRegistry(ctx context.Context) (*ent.Registry, error) {
	// Get default registry
	return self.base.DB.Registry.Query().
		Where(registry.IsDefault(true)).
		Order(
			ent.Desc(registry.FieldCreatedAt),
		).
		First(ctx)
}
