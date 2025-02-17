package repository

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/user"
)

func (r *Repository) GetOrCreateUser(ctx context.Context, email, username, subject string) (*ent.User, error) {
	// Try to find existing user
	user, err := r.DB.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}
		// Create new user if not found
		user, err = r.DB.User.
			Create().
			SetEmail(email).
			SetUsername(username).
			SetExternalID(subject).
			Save(ctx)

		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
