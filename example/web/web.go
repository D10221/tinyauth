package main

import (
	"net/http"
	"log"
	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
	"github.com/D10221/tinystore"
	"github.com/D10221/tinyauth/example/tinyapp"
	// "github.com/gorilla/mux"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var app = &tinyapp.TinyApp {Templates: "example/tinyapp"}

	config := config.NewConfig("")

	store:= &tinystore.SimpleStore{}
	app.Auth = tinyauth.NewTinyAuth(config, store)

	e := app.Auth.LoadConfig(app.MakePath("example/tinyapp/config.json"))

	if e != nil {
		panic(e)
	}

	if e = app.Auth.Config.Validate(); e != nil {
		panic(e)
	}

	e = store.LoadJson(app.MakePath("example/tinyapp/credentials.json"))
	if e != nil { panic(e)	}


	app.Auth.CredentialStore.ForEach(app.EncryptPassword)

	for _, c := range app.Auth.CredentialStore.All()[:]{
		log.Printf("Credentials: %v" , c)
	}

	log.Println(app.Auth.Config)
	log.Println(app.Auth.CredentialStore.All())
	//m:= mux.NewRouter()
	http.HandleFunc("/login", app.Login )
	http.HandleFunc("/authenticate", app.Authenticate)
	http.HandleFunc("/secret", app.Auth.RequireAuthentication(app.Secret))
	//Static

	path:= app.MakePath("/example/web/static")
	log.Printf("Serving static: %s", app.MakePath("/example/web/static"))
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	// http.Handle("/", m)
	address := ":8080"
	log.Printf("Serving : %v", address)
	http.ListenAndServe(address, nil)
}
