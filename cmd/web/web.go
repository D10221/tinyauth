package main

import (
	"net/http"
	"fmt"
	"github.com/D10221/tinyauth"
	"log"
	"path/filepath"
	"os"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "200 OK")
}

func main() {

	http.HandleFunc("/", app.Auth.RequireAuthentication(Hello))
	address := ":8080"
	log.Printf("ListenAndServe: %v", address)
	http.ListenAndServe(address, nil)
}
type TinyApp struct {
	Auth *tinyauth.TinyAuth
}
var app = &TinyApp {}

func init() {

	config:= &tinyauth.TinyAuthConfig{Secret: "", AuthorizationKey: "", BasicScheme: "" }
	app.Auth = tinyauth.NewTinyAuth(config)
	dir, e:= os.Getwd()
	if e!= nil {
		panic(e)
	}

	e= app.Auth.Config.LoadConfig(filepath.Join(dir,"cmd/web/config.json"))

	if e!= nil {
		panic(e)
	}
	if e = app.Auth.Config.Validate() ; e!=nil {
		panic(e)
	}
	if e= app.Auth.Config.Validate(); e!= nil{
		panic(e)
	}


	e = app.Auth.CredentialStore.LoadJson(filepath.Join(dir, "testdata/credentials.json"))
	if e!= nil {
		panic(e)
	}
}