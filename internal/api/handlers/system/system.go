package system_handler

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type HandlerGroup struct {
	srv *server.Server
}

func RegisterHandlers(server *server.Server, grp *huma.Group) {
	handlers := &HandlerGroup{
		srv: server,
	}

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-system-information",
			Summary:     "Get System Information",
			Description: "Get system information such as external load balancer IP for DNS configurations.",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetSystemInformation,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-system-settings",
			Summary:     "Update System Settings",
			Description: "Update system settings such as wild card domain, buildkit, etc.",
			Path:        "/settings/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateBuildkitSettings,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "check-dns-resolution",
			Summary:     "Check DNS Resolution",
			Description: "Check if the domain is pointing to the correct IP address.",
			Path:        "/dns/check",
			Method:      http.MethodGet,
		},
		handlers.CheckDNSResolution,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "check-for-updates",
			Summary:     "Check for Updates",
			Description: "Check for updates to the system.",
			Path:        "/upgrade/check",
			Method:      http.MethodGet,
		},
		handlers.CheckForUpdates,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "apply-update",
			Summary:     "Apply Update",
			Description: "Apply an update to the system.",
			Path:        "/upgrade/apply",
			Method:      http.MethodPost,
		},
		handlers.ApplyUpdate,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-upgrade-status",
			Summary:     "Get Upgrade Status",
			Description: "Get the status of the upgrade.",
			Path:        "/upgrade/status",
			Method:      http.MethodGet,
		},
		handlers.GetUpgradeStatus,
	)
}

func (self *HandlerGroup) handleErr(err error) error {
	if errors.Is(err, errdefs.ErrInvalidInput) {
		return huma.Error400BadRequest("invalid input", err)
	}
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound("entity not found", err)
	}
	log.Error("Unexpected system error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
