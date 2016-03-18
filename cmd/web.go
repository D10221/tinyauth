package main

import (
	"net/http"
	// "strings"
	// "encoding/base64"
	"fmt"
	"github.com/D10221/tinyauth/credentials"
)

var AuthorizationKey = "Authorization"

type Handler func (w http.ResponseWriter, r *http.Request);

func RequireAuthentication(handler Handler) Handler {
	return func (w http.ResponseWriter, r *http.Request){
		auth := r.Header.Get(AuthorizationKey)
		if(!credentials.ShouldDecode(auth).Authenticate()){
			http.Error(w, "unauthorized", 401)
			return
		}
		handler(w,r)
	}
}

func Hello(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "200 OK")
}

func main() {
	http.HandleFunc("/", RequireAuthentication(Hello))
	addr := ":8080"
	http.ListenAndServe(addr, nil)
}
