package user_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/server"
)

type MeResponse struct {
	Body struct {
		Data *ent.User `json:"data"`
	}
}

// Me handles GET /me
func (self *HandlerGroup) Me(ctx context.Context, _ *server.BaseAuthInput) (*MeResponse, error) {

	user, ok := ctx.Value("user").(*ent.User)
	if !ok {
		log.Error("Error getting user from context")
		return nil, huma.Error500InternalServerError("Unable to retrieve user")
	}

	resp := &MeResponse{}
	resp.Body.Data = user
	return resp, nil
}
