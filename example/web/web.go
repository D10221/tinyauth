package main

import (
	"net/http"
	"log"
	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
	"github.com/D10221/tinyauth/example/tinyapp"
	// "github.com/gorilla/mux"
	"github.com/D10221/tinyauth/store"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var app = &tinyapp.TinyApp {Templates: "cmd/tinyapp"}

	config := config.NewConfig("")

	app.Auth = tinyauth.NewTinyAuth(config)

	e := app.Auth.LoadConfig(app.MakePath("cmd/tinyapp/config.json"))

	if e != nil {
		panic(e)
	}

	if e = app.Auth.Config.Validate(); e != nil {
		panic(e)
	}

	e = app.Auth.CredentialStore.LoadJson(app.MakePath("cmd/tinyapp/credentials.json"))

	if e != nil {
		panic(e)
	}

	changePassword:= func(in *store.Credential) *store.Credential{
		password, e := app.Auth.Criptico.Encrypt(in.Password)
		if e!=nil {
			panic(e)
		}
		in.Password = password
		return in
	}

	app.Auth.CredentialStore.Update(changePassword)

	for _, c := range app.Auth.CredentialStore.All()[:]{
		log.Printf("Credentials: %v" , c)
	}

	log.Println(app.Auth.Config)
	log.Println(app.Auth.CredentialStore.All())
	//m:= mux.NewRouter()
	http.HandleFunc("/login", app.Login )
	http.HandleFunc("/authenticate", app.Authenticate)
	http.HandleFunc("/secret", app.Auth.RequireAuthentication(app.Secret))
	// http.Handle("/", m)
	address := ":8080"
	log.Printf("ListenAndServe: %v", address)
	http.ListenAndServe(address, nil)
}
