package oauth2server

import "net/http"

func (self *Oauth2Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	self.Srv.HandleTokenRequest(w, r)
}
