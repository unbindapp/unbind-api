package oauth2server

import (
	"encoding/json"
	"net/http"
)

func (self *Oauth2Server) HandleUserinfo(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid authorization header"})
		return
	}

	// The token is in the format "Bearer <token>"
	tokenString := authHeader[7:]

	// Get token info
	tokenInfo, err := self.Srv.Manager.LoadAccessToken(r.Context(), tokenString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid token"})
		return
	}

	// Get user info
	username := tokenInfo.GetUserID()
	u, err := self.Repository.User().GetByEmail(r.Context(), username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user info"})
		return
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"sub":   u.ID.String(),
		"email": u.Email,
	})
}
