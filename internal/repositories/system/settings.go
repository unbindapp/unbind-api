package system_repo

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *SystemRepository) GetSystemSettings(ctx context.Context, tx repository.TxInterface) (*ent.SystemSetting, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.SystemSetting.Query().
		First(ctx)
}

type SystemSettingUpdateInput struct {
	WildcardDomain *string
}

func (self *SystemRepository) UpdateSystemSettings(ctx context.Context, input *SystemSettingUpdateInput) (settings *ent.SystemSetting, err error) {
	if err := self.base.WithTx(ctx, func(tx repository.TxInterface) error {
		// Get system settings
		settings, err = self.GetSystemSettings(ctx, tx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if ent.IsNotFound(err) {
			// Create system settings
			settings, err = self.base.DB.SystemSetting.Create().Save(ctx)
			if err != nil {
				return err
			}
		}

		// Update system settings
		settings, err = tx.Client().SystemSetting.UpdateOneID(settings.ID).
			SetNillableWildcardBaseURL(input.WildcardDomain).
			Save(ctx)

		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return settings, nil
}
