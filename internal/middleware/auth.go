package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/log"
)

func (m *Middleware) Authenticate(ctx huma.Context, next func(huma.Context)) {
	// authHeader := ctx.Header("Authorization")
	// if authHeader == "" {
	// 	huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Authorization header required")
	// 	return
	// }

	// if !strings.HasPrefix(authHeader, "Bearer ") {
	// 	huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Authorization header must be a Bearer token")
	// 	return
	// }

	// bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
	bearerToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjE3YTk3YjBjYTc2OGQ1ZjgwM2M1MzI4ZjU0MDlhMmQzZTYyNzA3ZmEifQ.eyJpc3MiOiJodHRwczovL2RleC51bmJpbmQuYXBwIiwic3ViIjoiQ2lReE1qTTBOVFkzT0MweE1qTTBMVFUyTnpndE1USXpOQzAxTmpjNE1USXpORFUyTnpnU0JXeHZZMkZzIiwiYXVkIjoidW5iaW5kLWFwaSIsImV4cCI6MTc0MTI3MTA3MiwiaWF0IjoxNzQxMjcxMDEyLCJhdF9oYXNoIjoiaWp2ZzduZEZ6dWN5VDV6bjhYRGdhUSIsImVtYWlsIjoiYWRtaW5AdW5iaW5kLmFwcCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiYWRtaW4ifQ.nmHl7yAUa5QNMAOPF-uCDaQsiObS3IEPQiA0l0cORDDhBNEFumbrbBIqEfPTlvUHCKwEBsyBamd6PZbqUNHyvoMmMnhoAds4Om10vR3MnWBWr28Do3rGkrkiH4hJbqoJT1r163X4LFEuiPgV1m3K6QlKPFTZAT-IFQNMP6ki1mJQpKXN2gHzpYBaqpum24WgCoOILyKpISzNuUfuF_db1vR18bXKMQ13URgmG9mo-F0pxQHiKJCCvTQLGPvJ2eJTYqZOHDUfGsFCvpgPsMiE3KRewB5pWtEE0xpmERorETdAZeE_nmakjJE5bfeM5xziqdnK5sUvphZsOOdCArz38g"
	token, err := m.verifier.Verify(ctx.Context(), bearerToken)
	if err != nil {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Invalid token")
		return
	}

	var claims struct {
		Email    string `json:"email"`
		Username string `json:"preferred_username"`
		Subject  string `json:"sub"`
	}
	if err := token.Claims(&claims); err != nil {
		log.Errorf("Failed to parse claims: %v", err)
		huma.WriteErr(m.api, ctx, http.StatusInternalServerError, "Failed to parse claims")
		return
	}

	// Get or create user using Ent
	user, err := m.repository.GetOrCreateUser(ctx.Context(), claims.Email, claims.Username, claims.Subject)
	if err != nil {
		log.Errorf("Failed to process user: %v", err)
		huma.WriteErr(m.api, ctx, http.StatusInternalServerError, "Failed to process user")
		return
	}

	ctx = huma.WithValue(ctx, "user", user)
	next(ctx)
}
