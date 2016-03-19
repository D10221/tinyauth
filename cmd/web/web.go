package main

import (
	"net/http"
	"fmt"
	"github.com/D10221/tinyauth"
	"log"
	"github.com/D10221/tinyauth/credentials"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "200 OK")
}

func main() {
	http.HandleFunc("/", tinyauth.RequireAuthentication(Hello))
	address := ":8080"
	log.Printf("ListenAndServe: %v", address)
	http.ListenAndServe(address, nil)
}

func init() {
	credentials.Credentials.LoadJsonFromRelativePath("testdata/credentials.json")
}

