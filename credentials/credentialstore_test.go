package credentials_test

import (
	"testing"
	"github.com/D10221/tinyauth/credentials"
	"github.com/D10221/tinystore"
)

// Test_LoadJsonFile
func Test_LoadJsonFile(t *testing.T) {
	// new tinystore.Store , registering a adapter with its Own Name
	store := credentials.NewCredentialStore()
	e := tinystore.LoadJsonFile(store, "../testdata/credentials.json")
	if e != nil {
		t.Error(e)
		return
	}

	found, e := credentials.Find(store, credentials.UserNameEqualsFilter("admin"))

	if e != nil {
		t.Error(e)
		return
	}
	if value, ok := found.(*credentials.Credential); !ok {
		t.Error("Not *Credential?")
	} else {
		if value.Username != "admin" {
			t.Errorf("What we found?: &v", value)
		}
	}

}

// Test_LoadJson from byte[]
func Test_LoadJson(t *testing.T) {

	store := credentials.NewCredentialStore()
	bytes := []byte(`
	[
	  {
	    "Username": "admin","Password": "P@55w0rd!"
	  }
	]`)
	e := tinystore.LoadJson(store, bytes)
	if e != nil {
		t.Error(e)
		return
	}
	filter := func(name string) func(tinystore.StoreItem) bool {
		return func(item tinystore.StoreItem) bool {
			value, ok := item.(*credentials.Credential)
			if !ok {
				return false
			}
			return value.Username == name
		}
	}
	found, e := store.Find(filter("admin"))
	if e != nil {
		t.Error(e)
		return
	}
	if value, ok := found.(*credentials.Credential); !ok {
		t.Error("Not *Credential?")
	} else {
		if value.Username != "admin" {
			t.Errorf("What we found?: &v", value)
		}
	}

}