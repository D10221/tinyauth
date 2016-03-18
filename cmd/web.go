package main

import (
	"net/http"
	// "strings"
	// "encoding/base64"
	"fmt"
	"log"
	"github.com/D10221/tinyauth/credentials"
)



var AuthorizationKey = "Authorization"

func Authenticate(r *http.Request) (bool, error) {

	if r == nil {
		return false
	}

	// Confirm the request is sending Basic Authentication credentials.
	auth := r.Header.Get(AuthorizationKey)
	credentials, err := credentials.Decode(auth)
	if err!nil {
		log.Fatal(err)
	}
	return credentials.Authenticate()
}

func Handle(w http.ResponseWriter, r *http.Request) {
	ok, err := Authenticate(r)
	if !ok {
		log.Fatal(err)
		http.Error(w, "unauthorized", 401)
	}
	fmt.Fprint(w, "Ok")
}

func main() {
	http.HandleFunc("/", Handle)
	addr := ":8080"
	http.ListenAndServe(addr, nil)
}
