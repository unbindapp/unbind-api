package bootstrap_repo

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *BootstrapRepository) IsBootstrapped(ctx context.Context, tx repository.TxInterface) (bool, error) {
	// ! TODO - also check if users has any rows, but for now we are testing
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	// Check if bootstrap is already done
	_, err := db.Bootstrap.Query().First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
