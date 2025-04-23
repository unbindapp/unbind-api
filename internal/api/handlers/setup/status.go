package setup_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type SetupData struct {
	IsSetup bool `json:"is_setup"`
}

type SetupStatusResponse struct {
	Body struct {
		Data *SetupData `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) GetStatus(ctx context.Context, input *server.EmptyInput) (*SetupStatusResponse, error) {
	// Get bootstrapped
	bootstrapped, err := self.srv.Repository.Bootstrap().IsBootstrapped(ctx, nil)
	if err != nil {
		log.Error("Error checking if bootstrapped", "err", err)
		return nil, huma.Error500InternalServerError("Error checking if bootstrapped")
	}

	resp := &SetupStatusResponse{}
	resp.Body.Data = &SetupData{
		IsSetup: bootstrapped,
	}
	return resp, nil

}
