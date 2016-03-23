package tinyapp

import (
	"testing"
	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
	"log"
)

func Test_TinyApp(t *testing.T){

	app:=  &TinyApp{Auth: tinyauth.NewTinyAuth( config.NewConfig("0123456789ABCDEF"))}

	if app.Auth.Config.Secret == "" {
		t.Error("No Secret")
	}

	credential, e := app.Auth.NewCredential("user", "password")
	if e!=nil {
		t.Error(e)
		return;
	}
	if !credential.Valid() {
		t.Error("not valid credential")
	}
	if credential.Username != "user" {
		t.Error("Bad credential")
	}

	if credential.Password == "password" {
		t.Error("credential is not encrypted")
	}
}

func Test_TinyApp_Config(t *testing.T){

	app:=  &TinyApp{Auth: tinyauth.NewTinyAuth( &config.TinyAuthConfig{})}

	log.Printf("App Dir: %s", app.CurrentDir())

	path:= app.MakePath("config.json")

	log.Printf("App config: %s", path)

	e:= app.Auth.LoadConfig(path)

	if e!=nil {
		t.Error(e)
		return
	}

	if app.Auth.Config.Secret == "" {
		t.Error("No Secret")
		return
	}

	credential, e := app.Auth.NewCredential("user", "password")

	if e!=nil {
		t.Error(e)
		return
	}
	if !credential.Valid() {
		t.Error("not valid credential")
	}

	if credential.Username != "user" {
		t.Error("Bad credential")
	}

	if credential.Password == "password" {
		t.Error("credential is not encrypted")
	}


}

func Test_TinyApp_NoEncryption(t *testing.T){

	app:=  &TinyApp{Auth: tinyauth.NewTinyAuth( &config.TinyAuthConfig{})}

	if app.Auth.Config.Secret != "" {
		t.Error("Should not have Secret")
	}

	credential, e := app.Auth.NewCredential("user", "password")
	if e!= nil {
		t.Error(e)
		t.Fail()
		return
	}

	if !credential.Valid() {
		t.Error("not valid credential")
	}
	if credential.Username != "user" {
		t.Error("Bad credential user")
	}

	if credential.Password != "password" {
		t.Error("Bad credential password")
	}
	log.Println(credential)
}



