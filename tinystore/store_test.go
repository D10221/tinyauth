package tinystore_test

import (
	"testing"
	"github.com/D10221/tinyauth/tinystore"
)


func Test_TinyStoreFuncs(t *testing.T){
	var store tinystore.Store = &tinystore.SimpleStore{}
	c := &tinystore.Credential{"me", "1234"}
	e := store.Add(c)
	if e != nil { t.Error(e); return }

	t.Log("Where\n")
	if result,count := tinystore.Where(store, tinystore.UserNameEquals("me")) ; count != 1  || result[0].Username != "me"{
		t.Errorf("result: %v ,Count: %s", result, count)
	}
}


func Test_Where(t *testing.T){
	store:=  &tinystore.SimpleStore{}
	if e:= store.Add( &tinystore.Credential{"xyz", "1234"} ) ; e!= nil {
		t.Error(e)
		return
	}
	if e:= store.Add( &tinystore.Credential{"uuu", "1234"} ); e!=nil {
		t.Error(e)
		return
	}
	if len(store.All()[:])!= 2 {
		t.Error("Failed Add")
	}
	items, count := tinystore.Where(store, tinystore.UserNameEquals("xyz"))
	if count != 1 {
		t.Error("Bad Count")
		return
	}
	if items[0].Username != "xyz" {
		t.Error("Bad Item Returned")
	}

	notIems, count := tinystore.WhereNot(store, tinystore.UserNameEquals("xyz"))
	if count != 1 {
		t.Error("Bad Count")
		return
	}
	if notIems[0].Username!= "uuu" {
		t.Error("Bad Item Returned")
		return
	}
}

func Test_FindByKey_Load(t *testing.T){
	var store tinystore.Store = &tinystore.SimpleStore{}
	if e:= store.Load(&tinystore.Credential{"me", "1234"}, &tinystore.Credential{"el", "999"});e!=nil{
		t.Error(e)
		return
	}
	if c,e:= tinystore.FindByKey(store, "el"); e!=nil  || c.Username!= "el"{
		t.Error("Bad result")
	}
}

func Test_Store_ForEach(t *testing.T) {

	var store tinystore.Store = &tinystore.SimpleStore{}

	if e:= store.Load(&tinystore.Credential{"me", "1234"}, &tinystore.Credential{"el", "999"}); e!=nil {
		t.Error(e)
		return
	}

	changePasswords := func (c *tinystore.Credential) (*tinystore.Credential, error) {
		c.Password = "xxx" ;
		return c, nil
	}

	if e:= tinystore.ForEach(store, changePasswords, nil); e !=nil {
		t.Error(e);
		return
	}

	if x,e:= tinystore.FindByKey(store, "me"); e!=nil || x.Password != "xxx" {
		if e!=nil { t.Error(e) ; return }
		t.Errorf("Not ok => x: %v", x)
	}


}
