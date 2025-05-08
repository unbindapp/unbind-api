package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type ListPVCInput struct {
	server.BaseAuthInput
	models.ListPVCInput
}

type ListPVCResponse struct {
	Body struct {
		Data []k8s.PVCInfo `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) ListPVCs(ctx context.Context, input *ListPVCInput) (*ListPVCResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	pvcs, err := self.srv.StorageService.ListPVCs(ctx, user.ID, bearerToken, &input.ListPVCInput)
	if err != nil {
		return nil, self.handleErr(err)
	}
	if len(pvcs) == 0 {
		pvcs = []k8s.PVCInfo{}
	}

	response := &ListPVCResponse{}
	response.Body.Data = pvcs
	return response, nil
}

// * get by ID
type GetPVCInput struct {
	server.BaseAuthInput
	models.GetPVCInput
}

type GetPVCResponse struct {
	Body struct {
		Data *k8s.PVCInfo `json:"data"`
	}
}

func (self *HandlerGroup) GetPVC(ctx context.Context, input *GetPVCInput) (*GetPVCResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	pvc, err := self.srv.StorageService.GetPVC(ctx, user.ID, bearerToken, &input.GetPVCInput)
	if err != nil {
		return nil, self.handleErr(err)
	}

	response := &GetPVCResponse{}
	response.Body.Data = pvc
	return response, nil
}
