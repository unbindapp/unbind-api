package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

func (self *Middleware) CheckBootstrapped(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader != "" {
		// Don't check bootstrap if we have an auth header
		next(ctx)
		return
	}

	bs, err := self.repository.Ent().Bootstrap.Query().First(ctx.Context())
	if err != nil && !ent.IsNotFound(err) {
		log.Error("Failed to check bootstrap status: %v", err)
		huma.WriteErr(self.api, ctx, http.StatusInternalServerError, "Server failed to execute query")
		return
	}

	if bs == nil || ent.IsNotFound(err) {
		cookie := http.Cookie{
			Name:  "needs-bootstrap",
			Value: "true",
		}
		ctx.AppendHeader("Set-Cookie", cookie.String())
	}

	next(ctx)
}
