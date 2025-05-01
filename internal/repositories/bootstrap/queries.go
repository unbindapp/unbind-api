package bootstrap_repo

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *BootstrapRepository) IsBootstrapped(ctx context.Context, tx repository.TxInterface) (userExists bool, isBootstrapped bool, err error) {
	// ! TODO - also check if users has any rows, but for now we are testing
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	// Check if user exists
	userCount, err := db.User.Query().
		Count(ctx)
	if err != nil {
		return false, false, err
	}

	userExists = userCount > 0

	// Check if bootstrap is already done
	_, err = db.Bootstrap.Query().First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return userExists, false, nil
		}
		return userExists, false, err
	}
	return userExists, true, nil
}
