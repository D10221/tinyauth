package tinystore_test

import (
	"testing"
	"fmt"
	"os"
	"path/filepath"
	"github.com/D10221/tinyauth/tinystore"
)


func reverse(s string) string {
	var result []byte
	for i := len([]byte(s))-1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return string(result)
}

func Test_NewLocalStore(t *testing.T) {

	store:= & tinystore.SimpleStore{}

	t.Log("Add")
	e := store.Add(&tinystore.Credential{"me", "1234"})
	if e != nil {
		t.Error(e); return
	}

	t.Log("FindBy ByName")
	if x, e := store.Find(tinystore.UserNameEquals("me")); e != nil || x.Username != "me" {
		t.Error(e)
		return
	}

	reversePassword := func(c *tinystore.Credential) (*tinystore.Credential, error) {
		c.Password = reverse(c.Password)
		return c, nil
	}

	t.Log("UpdateAll")
	e = store.ForEach(reversePassword)
	if e!=nil {
		t.Error(e)
		return
	}


	t.Log("FindBy")
	if x,e:= store.Find(tinystore.UserNameEquals("me")); e!=nil || x.Password!= "4321" {
		if e!=nil {
			t.Error(e)
			return
		}
		t.Error("UpdateAll Fail")
		return
	}


	t.Log("UpdateWhere")
	if ex := store.Add(&tinystore.Credential{"you", "1234"}) ; ex != nil {
		t.Error(e); return
	}
	if ex:= store.ForEachWhere(tinystore.UserNameEquals("you"), reversePassword) ; ex!=nil {
		t.Error(ex)
		return
	}

	t.Log("FindBy")
	if x,e:= store.Find(tinystore.UserNameEquals("you")); e!=nil || x.Password!= "4321" {
		if e!=nil {
			t.Error(e)
			return
		}
		t.Error("UpdateWhere Fail")
		return
	}


	t.Log("Remove")
	if ex:= store.Remove(&tinystore.Credential{"me", "1234"}) ; ex!=nil {
		t.Error(ex)
		return
	}

	t.Log("FindBy")
	if x, ex:= store.Find(tinystore.UserNameEquals("me")); ex!= tinystore.ErrNotFound {
		t.Log(x)
		if ex!=nil {
			t.Error(ex)
			return
		}
		if x.Valid() {
			t.Error("Wtf")
		}
		t.Error("Not FOund?")
	}

	t.Log("RemoveWhere")

	if found , ex := store.Find(tinystore.UserNameEquals("you")) ; ex != nil || found.Username != "you"{
		t.Error(ex)
		return
	}
	if ex:= store.RemoveWhere(tinystore.UserNameEquals("you")) ; ex!=nil {
		t.Error(ex)
		return
	}
	if found , ex := store.Find(tinystore.UserNameEquals("you")) ; ex != tinystore.ErrNotFound {
		t.Log(found)
		t.Error(ex)
		return
	}

	t.Log("Length")
	if l,e := store.Length(); e!= nil || l != 0 {
		if e!=nil {
			t.Error(e)
			return
		}
		t.Error("Failed count")
		return
	}

	for i:= 0; i <= 99 ; i++ {
		name:= fmt.Sprintf("%s", i)
		store.Add(&tinystore.Credential{name, "1234"})
	}
	if l,e:= store.Length() ; e!=nil || l != 100 {
		if e!=nil {
			t.Error(e)
			return
		}
		t.Errorf("Bad Length: %s", l )
	}

	t.Log("Clear")

	store.Clear()

	if l,e:= store.Length() ; e!=nil || l != 0 {
		if e!=nil {
			t.Error(e)
			return
		}
		t.Errorf("Bad Length: %s", l )
	}

}


func Test_Store(t *testing.T) {

	store := &tinystore.SimpleStore{}

	store.Add(&tinystore.Credential{"admin", "password"})

	if store.All()[0].Username != "admin" {
		t.Error("Bad Store")
	}

	credential, e := store.Find(tinystore.UserNameEquals("admin"))
	if e != nil {
		t.Error(e); return
	}

	if credential.Username != "admin" || credential.Password != "password" {
		t.Error("Bad store")
	}

	store.Clear()
	store.Add(&tinystore.Credential{Username: "me", Password:"1234"})

	user, e := store.Find(tinystore.UserNameEquals("me"))

	if e != nil {
		t.Error(e); return
	}

	if user.Username != "me" {
		t.Error("Not found")
	}

	if tinystore.Length(store) != 1 {
		t.Error("Wtf")
	}

}

func Test_Store_json_load(t *testing.T) {

	store := &tinystore.SimpleStore{}

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

	credential, e := store.Find(tinystore.UserNameEquals("admin"))
	if e != nil {
		t.Error(e); return
	}

	if credential.Username != "admin" || credential.Password != "P@55w0rd!" {
		t.Error("Bad store")
	}

	store.Clear()
	store.Add(&tinystore.Credential{Username: "me", Password:"1234"})

	user, e := store.Find(tinystore.UserNameEquals("me"))
	if e != nil {
		t.Error(e); return
	}

	if user.Username != "me" {
		t.Error("Not found")
	}

	if tinystore.Length(store) != 1 {
		t.Error("Wtf")
	}
}

func Test_Add(t *testing.T) {
	store := tinystore.SimpleStore{}
	e := store.Add(&tinystore.Credential{"me", "1234"})
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
	store := &tinystore.SimpleStore{}
	e := store.Add(&tinystore.Credential{"me", "1234"})
	if e != nil {
		t.Error(e)
	}
	e = store.Add(&tinystore.Credential{"me", "1234"})
	if e != tinystore.ErrAlreadyExists {
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
	store := tinystore.SimpleStore{}
	e := store.Add(&tinystore.Credential{"", ""})
	if e != tinystore.ErrInvalidCredential {
		t.Error("Should return InvalidCredential")
	}
	e = nil
	e = store.Add(&tinystore.Credential{})
	if e != tinystore.ErrInvalidCredential {
		t.Error("Should return InvalidCredential")
	}
	if len(store.All()) != 0 {
		t.Error("Failed len")
	}
}

func Test_Remove(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	e := store.Remove(&tinystore.Credential{})

	if e != tinystore.ErrInvalidCredential {
		t.Error("Should be invalid")
		return
	}
	e = nil
	e = store.Remove(&tinystore.Credential{"me", "1234"})
	if e != tinystore.ErrNotFound {
		t.Error("Should be NotFound")
		return
	}
	e = nil
	e = store.Add(&tinystore.Credential{"me", "1234"})
	if e != nil {
		t.Error(e); return
	}
	e = store.Remove(&tinystore.Credential{"me", "1234"})
	if e != nil {
		t.Errorf("Shouldn't error %s", e.Error())
		return
	}
	if tinystore.Length(store) != 0 {
		t.Error("len fail")
	}
}

func Test_Remove_Add(t *testing.T) {

	var store tinystore.Store  = &tinystore.SimpleStore{}

	e := store.Add(&tinystore.Credential{"me", "1234"})
	// password is not Checked
	e = store.Remove(&tinystore.Credential{"me", "1111"})
	if e != nil {
		t.Error(e); return
	}
	e = store.Add(&tinystore.Credential{"me", "4321"})
	if e != nil {
		t.Error(e); return
	}
	if len(store.All()) != 1 {
		t.Error("Len Fail")
	}
	if store.All()[0].Username != "me" && store.All()[0].Password != "4321" {
		t.Error("Wrong values")
	}
	if u, e := tinystore.FindByKey(store,"me") ; e != nil && !u.Valid() || u.Username != "me" || u.Password != "4321" {
		t.Error("Wtf")
	}
}

func Test_Update(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	e := store.Add( &tinystore.Credential{"me", "1234"} )

	if e != nil {
		t.Error(e)
	}

	changePassword := func(in *tinystore.Credential) (*tinystore.Credential, error) {
		in.Password = "abcd"
		return in, nil
	}

	e = store.ForEach(changePassword)
	if e != nil {
		t.Error(e)
	}
	if len(store.All()) != 1 || store.All()[0].Password != "abcd" {
		t.Error("Doesn't work")
	}
}

func Test_UpdateWhere(t *testing.T) {

	store := tinystore.SimpleStore{}

	e := store.Add(&tinystore.Credential{"me", "1234"})

	if e != nil {
		t.Error(e)
	}

	changePassword := func(in *tinystore.Credential) (*tinystore.Credential, error) {
		if in.Password == "" {
			return nil, tinystore.ErrInvalidCredential
		}
		in.Password = "abcd"
		return in, nil
	}

	e = store.ForEachWhere(tinystore.UserNameEquals("me"), changePassword)
	if e != nil {
		t.Error(e)
	}
	if len(store.All()) != 1 || store.All()[0].Password != "abcd" {
		t.Error("Doesn't work")
	}
}

func Test_FindBy(t *testing.T) {

	store := tinystore.SimpleStore{}

	e := store.Add(&tinystore.Credential{"me", "1234"})

	if e != nil {
		t.Error(e); return
	}

	found, e := store.Find(tinystore.UserNameEquals("me"))

	if e != nil || ! found.Valid() {
		t.Error(e)
	}

	if found.Username != "me" {
		t.Error("Wtf")
	}
}

func Test_RemoveWhere(t *testing.T) {

	store := &tinystore.SimpleStore{}

	e := store.Add(&tinystore.Credential{"me", "1234"})
	if e != nil {
		t.Error(e); return
	}

	e = store.RemoveWhere(tinystore.UserNameEquals("me"))
	if e != nil {
		t.Error(e)
	}
	if found, e := store.Find(tinystore.UserNameEquals("me")); e != nil &&  e != tinystore.ErrNotFound || found.Valid() {
		t.Error("Shouldn't be found")
	}

}

