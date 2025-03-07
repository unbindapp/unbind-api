package oauth2server

import (
	"crypto/rsa"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database/repository"
)

type Oauth2Server struct {
	Cfg        *config.Config
	Repository *repository.Repository
	Srv        *server.Server
	PrivateKey *rsa.PrivateKey
	Kid        string
}
