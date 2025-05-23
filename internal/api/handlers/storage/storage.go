package storage_handler

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
			OperationID: "test-s3-access",
			Summary:     "Test S3 Access",
			Description: "Test S3 access with the provided credentials.",
			Path:        "/s3/test",
			Method:      http.MethodPost,
		},
		handlers.TestS3Access,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-s3-source",
			Summary:     "Create S3 Source",
			Description: "Create an S3 source to be used for backups, etc.",
			Path:        "/s3/create",
			Method:      http.MethodPost,
		},
		handlers.CreateS3,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-s3-source",
			Summary:     "Update S3 Source",
			Description: "Update an S3 source.",
			Path:        "/s3/update",
			Method:      http.MethodPost,
		},
		handlers.UpdateS3Source,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-s3-source-by-id",
			Summary:     "Get S3 Source by ID",
			Description: "Get S3 source ID.",
			Path:        "/s3/get",
			Method:      http.MethodGet,
		},
		handlers.GetS3SourceByID,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-s3-sources",
			Summary:     "List S3 Sources",
			Description: "List all S3 sources for a team.",
			Path:        "/s3/list",
			Method:      http.MethodGet,
		},
		handlers.ListS3Source,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-s3-source",
			Summary:     "Delete S3 Source",
			Description: "Delete an S3 source.",
			Path:        "/s3/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteS3Source,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-pvc",
			Summary:     "List PVCs",
			Description: "List all PVCs for a team, project, or environment.",
			Path:        "/pvc/list",
			Method:      http.MethodGet,
		},
		handlers.ListPVCs,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-pvc",
			Summary:     "Get PVC",
			Description: "Get a PVC by name.",
			Path:        "/pvc/get",
			Method:      http.MethodGet,
		},
		handlers.GetPVC,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-pvc",
			Summary:     "Create PVC",
			Description: "Create a PVC.",
			Path:        "/pvc/create",
			Method:      http.MethodPost,
		},
		handlers.CreatePVC,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-pvc",
			Summary:     "Update PVC",
			Description: "Update a PVC.",
			Path:        "/pvc/update",
			Method:      http.MethodPut,
		},
		handlers.UpdatePVC,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-pvc",
			Summary:     "Delete PVC",
			Description: "Delete a PVC.",
			Path:        "/pvc/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeletePVC,
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
	log.Error("Unexpected storage error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
