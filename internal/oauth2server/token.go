package oauth2server

import "net/http"

func (self *Oauth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	_ = self.Srv.HandleTokenRequest(w, r)
}
