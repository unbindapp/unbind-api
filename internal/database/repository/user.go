package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/user"
	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) GetOrCreateUser(ctx context.Context, email string) (*ent.User, error) {
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
			Save(ctx)

		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.DB.User.Query().Where(user.EmailEQ(email)).Only(ctx)
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidUserInput = errors.New("invalid user input")
)

// AuthenticateUser verifies a user's credentials and returns the user if successful
func (r *Repository) AuthenticateUser(ctx context.Context, email, password string) (*ent.User, error) {
	if email == "" || password == "" {
		return nil, ErrInvalidUserInput
	}

	// Find the user by email
	user, err := r.DB.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("querying user: %w", err)
	}

	// Verify password using bcrypt
	if err := verifyPassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// verifyPassword checks if the provided password matches the stored hash
func verifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
