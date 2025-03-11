package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
)

func (self *Server) GetUserFromContext(ctx context.Context) (user *ent.User, found bool) {
	user, found = ctx.Value("user").(*ent.User)
	return user, found
}

type MeResponse struct {
	Body struct {
		User *ent.User `json:"user"`
	}
}

// Me handles GET /me
func (self *Server) Me(ctx context.Context, _ *EmptyInput) (*MeResponse, error) {

	user, ok := ctx.Value("user").(*ent.User)
	if !ok {
		log.Error("Error getting user from context")
		return nil, huma.Error500InternalServerError("Unable to retrieve user")
	}

	resp := &MeResponse{}
	resp.Body.User = user
	return resp, nil
}
