package store

import (
	"testing"
	"path/filepath"
	"os"
)

func Test_Store(t *testing.T) {

	store := &SimpleCredentialStore{}

	store.Load(&Credential{"admin", "password"})

	if store.All()[0].Username != "admin" {
		t.Error("Bad Store")
	}

	credential := store.FindUser("admin")

	if credential.Username != "admin" || credential.Password != "password" {
		t.Error("Bad store")
	}

	store.Load(&Credential{Username: "me", Password:"1234"})

	user := store.FindUser("me")

	if user.Username != "me" {
		t.Error("Not found")
	}

	if len(store.all) != 1 {
		t.Error("Wtf")
	}

}

func Test_Store_json_load(t *testing.T) {

	store := &SimpleCredentialStore{}

	dir, e := os.Getwd()
	if e != nil {
		t.Error(e)
	}
	path := filepath.Join(dir, "../testdata/credentials.json")
	_, e = os.Stat(path)
	if e != nil {
		t.Error(e)
		return
	}

	e = store.LoadJson(path)

	if e != nil {
		t.Error(e)
	}

	if store.All()[0].Username != "admin" {
		t.Error("Bad Store")
	}

	credential := store.FindUser("admin")

	if credential.Username != "admin" || credential.Password != "P@55w0rd!" {
		t.Error("Bad store")
	}

	store.Load(&Credential{Username: "me", Password:"1234"})

	user := store.FindUser("me")

	if user.Username != "me" {
		t.Error("Not found")
	}

	if len(store.all) != 1 {
		t.Error("Wtf")
	}
}

func Test_Add(t *testing.T) {
	store := SimpleCredentialStore{}
	e := store.Add(&Credential{"me", "1234"})
	if e != nil {
		t.Error(e)
	}
	if len(store.All()) != 1 {
		t.Error("Failed len")
	}
	if store.All()[0].Username != "me" {
		t.Error("Fail add")
	}
}

func Test_Add_Existing(t *testing.T) {
	store := SimpleCredentialStore{}
	e := store.Add(&Credential{"me", "1234"})
	if e != nil {
		t.Error(e)
	}
	e = store.Add(&Credential{"me", "1234"})
	if e != AlreadyExists {
		t.Error("Should return AlreadyExists")
	}
	if len(store.All()) != 1 {
		t.Error("Failed len")
	}
	if store.All()[0].Username != "me" {
		t.Error("Fail add")
	}
}

func Test_Add_Empty(t *testing.T) {
	store := SimpleCredentialStore{}
	e := store.Add(&Credential{"", ""})
	if e != InvalidCredential {
		t.Error("Should return InvalidCredential")
	}
	e = nil
	e = store.Add(&Credential{})
	if e != InvalidCredential {
		t.Error("Should return InvalidCredential")
	}
	if len(store.All()) != 0 {
		t.Error("Failed len")
	}
}

func Test_Remove(t *testing.T) {

	store := SimpleCredentialStore{}

	e := store.Remove(&Credential{})

	if e != InvalidCredential {
		t.Error("Should be invalid")
		return
	}
	e = nil
	e = store.Remove(&Credential{"me", "1234"})
	if e != NotFound {
		t.Error("Should be NotFound")
		return
	}
	e = nil
	e = store.Add(&Credential{"me", "1234"})
	if e != nil {
		t.Error(e); return
	}
	e = store.Remove(&Credential{"me", "1234"})
	if e != nil {
		t.Errorf("Shouldn't error %s", e.Error())
		return
	}
	if len(store.all) != 0 {
		t.Error("len fail")
	}
}

func Test_Remove_Add(t *testing.T) {
	store := SimpleCredentialStore{}
	e:= store.Add(&Credential{"me", "1234"})
	// paassword is not Checked
	e = store.Remove(&Credential{"me" , "1111"})
	if e!= nil { t.Error(e); return }
	e = store.Add(&Credential{"me", "4321"})
	if e!= nil { 	t.Error(e); return }
	if len(store.All()) != 1 {
		t.Error("Len Fail")
	}
	if store.All()[0].Username != "me" && store.All()[0].Password != "4321" {
		t.Error("Wrong values")
	}
	if u:= store.FindUser("me") ; !u.Valid() || u.Username!= "me" ||  u.Password!= "4321" {
		t.Error("Wtf")
	}
}

func Test_Update(t *testing.T){

	store:= SimpleCredentialStore{}

	e:= store.Add(&Credential{"me", "1234"})

	if e!= nil { t.Error(e) }

	changePassword:= func(in *Credential) *Credential{
		in.Password = "abcd"
		return in
	}

	e = store.Update(changePassword)
	if e!=nil {
		t.Error(e)
	}
	if  len(store.All()) != 1 || store.All()[0].Password!= "abcd" {
		t.Error("Doesn't work")
	}
}
