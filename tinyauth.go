package tinyauth

import (
	"net/http"
	"github.com/D10221/tinyauth/credentials"
	"github.com/D10221/tinyauth/config"
)

type Handler func (w http.ResponseWriter, r *http.Request);

func RequireAuthentication(handler Handler) Handler {
	return func (w http.ResponseWriter, r *http.Request){
		auth := r.Header.Get(config.Current.AuthorizationKey)
		if(!credentials.ShouldDecode(auth).Authenticate()){
			http.Error(w, "unauthorized", 401)
			return
		}
		handler(w,r)
	}
}
