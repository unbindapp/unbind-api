package oauth2server

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (self *Oauth2Server) PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	// Find the user
	u, err := self.Repository.GetUserByEmail(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	return username, nil
}
