package setup_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type SetupData struct {
	Bootstrapped bool `json:"bootstrapped"`
}

type SetupStatusResponse struct {
	Body struct {
		Data *SetupData `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) GetStatus(ctx context.Context, input *server.EmptyInput) (*SetupStatusResponse, error) {
	if self.setupDone {
		resp := &SetupStatusResponse{}
		resp.Body.Data = &SetupData{
			Bootstrapped: true,
		}
		return resp, nil
	}

	// Get bootstrapped
	bootstrapped, err := self.srv.Repository.Bootstrap().IsBootstrapped(ctx, nil)
	if err != nil {
		log.Error("Error checking if bootstrapped", "err", err)
		return nil, huma.Error500InternalServerError("Error checking if bootstrapped")
	}

	if bootstrapped {
		self.setupDone = true
	}

	resp := &SetupStatusResponse{}
	resp.Body.Data = &SetupData{
		Bootstrapped: bootstrapped,
	}
	return resp, nil
}
