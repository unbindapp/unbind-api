package oauth2server

import (
	"encoding/json"
	"net/http"

	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// WellKnownConfig holds the fields you want to return in /.well-known/openid-configuration
type WellKnownConfig struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint,omitempty"`
	JWKSURI                          string   `json:"jwks_uri,omitempty"`
	ResponseTypesSupported           []string `json:"response_types_supported,omitempty"`
	SubjectTypesSupported            []string `json:"subject_types_supported,omitempty"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported,omitempty"`
	ScopesSupported                  []string `json:"scopes_supported,omitempty"`
}

// HandleWellKnown serves the OIDC Discovery document at /.well-known/openid-configuration
func (self *Oauth2Server) HandleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	authorizeEndpoint, _ := utils.JoinURLPaths(self.Cfg.ExternalOauth2URL, "authorize")
	tokenEndpoint, _ := utils.JoinURLPaths(self.Cfg.ExternalOauth2URL, "token")
	userinfoEndpoint, _ := utils.JoinURLPaths(self.Cfg.ExternalOauth2URL, "userinfo")
	jwksEndpoint, _ := utils.JoinURLPaths(self.Cfg.ExternalOauth2URL, ".well-known", "jwks.json")

	// You might construct these URLs dynamically based on your server’s configuration
	wellKnown := WellKnownConfig{
		Issuer:                self.Cfg.ExternalOauth2URL,
		AuthorizationEndpoint: authorizeEndpoint,
		TokenEndpoint:         tokenEndpoint,
		UserinfoEndpoint:      userinfoEndpoint,
		JWKSURI:               jwksEndpoint,
		ResponseTypesSupported: []string{
			"code",
		},
		SubjectTypesSupported: []string{
			"public",
		},
		IDTokenSigningAlgValuesSupported: []string{
			"RS256",
		},
		ScopesSupported: []string{
			"openid",
			"profile",
			"email",
			"offline_access",
			"groups",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wellKnown); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
