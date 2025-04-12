package metrics_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type GetMetricsInput struct {
	server.BaseAuthInput
	models.MetricsQueryInput
}

type GetMetricsResponse struct {
	Body struct {
		Data *models.MetricsResult `json:"data"`
	}
}

func (self *HandlerGroup) GetMetrics(ctx context.Context, input *GetMetricsInput) (*GetMetricsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	metrics, err := self.srv.MetricsService.GetMetrics(ctx, user.ID, &input.MetricsQueryInput)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &GetMetricsResponse{}
	resp.Body.Data = metrics
	return resp, nil
}

// Node level

type GetNodeMetricsInput struct {
	server.BaseAuthInput
	models.NodeMetricsQueryInput
}

type GetNodeMetricsResponse struct {
	Body struct {
		Data *models.NodeMetricsResult `json:"data"`
	}
}

func (self *HandlerGroup) GetNodeMetrics(ctx context.Context, input *GetNodeMetricsInput) (*GetNodeMetricsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	metrics, err := self.srv.MetricsService.GetNodeMetrics(ctx, user.ID, &input.NodeMetricsQueryInput)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &GetNodeMetricsResponse{}
	resp.Body.Data = metrics
	return resp, nil
}
