package service_handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/pkg/databases"
)

type ListDatabasesResponse struct {
	Body struct {
		Data []*databases.DatabaseList `json:"data" nullable:"false"`
	}
}

// ListDatabases handles GET /databases/list
func (self *HandlerGroup) ListDatabases(ctx context.Context, input *server.BaseAuthInput) (*ListDatabasesResponse, error) {

	dbList, err := self.srv.DatabaseProvider.ListDatabases(ctx, self.srv.Cfg.UnbindServiceDefVersion)
	if err != nil {
		log.Errorf("failed to list databases: %v", err)
		return nil, huma.Error500InternalServerError("An unknown error occured")
	}

	return &ListDatabasesResponse{
		Body: struct {
			Data []*databases.DatabaseList `json:"data" nullable:"false"`
		}{
			Data: []*databases.DatabaseList{dbList},
		},
	}, nil
}

type GetDatabaseSpecInput struct {
	server.BaseAuthInput
	Type    string `query:"type" required:"true" description:"Type of the database resource (e.g. postgres)"`
	Version string `query:"version" required:"false" description:"Version of the custom services release"`
}

type GetDatabaseResponse struct {
	Body struct {
		Data *databases.Definition `json:"data" nullable:"false"`
	}
}

// ListDatabases handles GET /databases/get
func (self *HandlerGroup) GetDatabaseDefinition(ctx context.Context, input *GetDatabaseSpecInput) (*GetDatabaseResponse, error) {
	version := self.srv.Cfg.UnbindServiceDefVersion
	if input.Version != "" {
		version = input.Version
	}

	template, err := self.srv.DatabaseProvider.FetchDatabaseDefinition(
		ctx,
		version,
		input.Type,
	)
	if err != nil {
		if errors.Is(err, databases.ErrDatabaseNotFound) {
			return nil, huma.Error404NotFound("Database not found")
		}
		log.Errorf("failed to get databases: %v", err)
		return nil, huma.Error500InternalServerError("An unknown error occured")
	}

	response := &GetDatabaseResponse{}
	response.Body.Data = template
	return response, nil
}
