package oauth2server

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
)

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"` // "RSA"
	N   string `json:"n"`   // Base64url-encoded modulus
	E   string `json:"e"`   // Base64url-encoded exponent
	Alg string `json:"alg"` // "RS256"
	Use string `json:"use"` // "sig"
	Kid string `json:"kid"` // Key ID
}

// HandleJWKS returns your server's public key in JWKS format
func (self *Oauth2Server) HandleJWKS(w http.ResponseWriter, r *http.Request) {
	if self.PrivateKey == nil {
		http.Error(w, "no private key configured", http.StatusInternalServerError)
		return
	}

	pubKey := self.PrivateKey.PublicKey

	// Convert modulus (N) and exponent (E) to base64-url strings
	n := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())

	eBytes := big.NewInt(int64(pubKey.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	// The "kid" should match what you used in NewJWTAccessGenerate
	jwkResp := jwk{
		Kty: "RSA",
		N:   n,
		E:   e,
		Alg: "RS256",
		Use: "sig",
		Kid: self.Kid,
	}

	resp := jwksResponse{Keys: []jwk{jwkResp}}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
