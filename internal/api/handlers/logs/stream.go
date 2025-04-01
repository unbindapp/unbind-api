package logs_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// LogEvent represents a log line event sent via SSE
type LogEvent struct {
	PodName   string    `json:"podName"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message"`
}

// Parameters for querying logs
type GetLogInput struct {
	server.BaseAuthInput
	models.LogStreamInput
}

func (self *HandlerGroup) GetLogsfunc(ctx context.Context, input *GetLogInput, send sse.Sender) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		send.Data(
			loki.LogsError{
				Code:    http.StatusUnauthorized,
				Message: "Unable to retrieve user",
			},
		)
	}

	if err := self.srv.LogService.StreamLogs(ctx, user.ID, &input.LogStreamInput, send); err != nil {
		self.handleSSEErr(err, send)
	}
}
