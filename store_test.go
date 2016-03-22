package tinyauth

import (
	"testing"
	"path/filepath"
	"os"
)

func Test_Store(t *testing.T){

	simpleStore := &SimpleCredentialStore{}

	dir, e := os.Getwd()
	if e != nil {
		t.Error(e)
	}

	simpleStore.LoadJson(filepath.Join(dir, "testdata/credentials.json"))

	var store CredentialStore = simpleStore

	if store.All()[0].Username != "admin" {
		t.Error("Bad Store")
	}

	credential, e := store.FindUser("admin")
	if e!=nil{
		t.Error(e)
	}

	if credential.Username!= "admin" || credential.Password != "P@55w0rd!"{
		t.Error("Bad store")
	}

	simpleStore.Load(Credential{Username: "me", Password:"1234"})
	user,e:= simpleStore.FindUser("me")
	if e!=nil {
		t.Error(e)
	}
	if user.Username != "me" {
		t.Error("Not found")
	}

	if len(simpleStore.all) != 1 {
		t.Error("Wtf")
	}

}
