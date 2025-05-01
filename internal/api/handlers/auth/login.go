package auth_handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/oauth2server"
)

type LoginForm struct {
	Username      string `schema:"username"  json:"username"  form:"username"  minLength:"1"`
	Password      string `schema:"password"  json:"password"  form:"password"  minLength:"1"`
	ClientID      string `schema:"client_id" json:"client_id" form:"client_id" minLength:"1"`
	RedirectURI   string `schema:"redirect_uri" json:"redirect_uri" form:"redirect_uri" format:"uri"`
	ResponseType  string `schema:"response_type" json:"response_type" form:"response_type" default:"code"`
	State         string `schema:"state" json:"state" form:"state"`
	Scope         string `schema:"scope" json:"scope" form:"scope"`
	PageKey       string `schema:"page_key" json:"page_key" form:"page_key"`
	InitiatingURL string `schema:"initiating_url" json:"initiating_url" form:"initiating_url"`
}

type LoginSubmitInput struct {
	Body LoginForm
}

type LoginSubmitResponse struct {
	Status   int
	Location string `header:"Location"`
}

func (h *HandlerGroup) LoginSubmit(
	ctx context.Context,
	input *LoginSubmitInput,
) (*LoginSubmitResponse, error) {
	user, err := h.srv.Repository.User().
		Authenticate(ctx, input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	authorizeURL, err := oauth2server.BuildOauthRedirect(
		h.srv.Cfg, oauth2server.RedirectAuthorize,
		map[string]string{
			"client_id":      input.Body.ClientID,
			"redirect_uri":   input.Body.RedirectURI,
			"response_type":  input.Body.ResponseType,
			"state":          input.Body.State,
			"scope":          input.Body.Scope,
			"user_id":        user.Email,
			"initiating_url": input.Body.InitiatingURL,
		})
	if err != nil {
		return nil, huma.Error500InternalServerError(
			"could not create authorize redirect", err)
	}

	return &LoginSubmitResponse{
		Status:   http.StatusFound,
		Location: authorizeURL,
	}, nil
}
