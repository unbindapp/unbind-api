package setup_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type CreateUserInput struct {
	Body struct {
		Email    string `json:"email" required:"true" format:"email"`
		Password string `json:"password" required:"true" minLength:"6"`
	}
}

type UserData struct {
	Email string `json:"email" required:"true" format:"email"`
}

type CreateUserResponse struct {
	Body struct {
		Data *UserData `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserResponse, error) {
	// Create user
	user, err := self.srv.Repository.Bootstrap().CreateUser(
		ctx,
		input.Body.Email,
		input.Body.Password,
	)
	if err != nil {
		log.Error("Error creating user", "err", err)
		return nil, huma.Error500InternalServerError("Error creating user")
	}

	resp := &CreateUserResponse{}
	resp.Body.Data = &UserData{
		Email: user.Email,
	}
	return resp, nil

}
