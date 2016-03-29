package tinyauth

import (
	"github.com/D10221/tinystore"
	"github.com/D10221/tinyauth/credentials"
)

var convert = func(item map[string]interface{}) tinystore.StoreItem {

	newItem:= &credentials.Credential{}

	if keyValue, keyExists := item["Username"]; keyExists  {
		if value, ok := keyValue.(string); ok {
			newItem.Username = value
		}
	}

	if keyValue, keyExists := item["Password"]; keyExists  {
		if value, ok := keyValue.(string); ok {
			newItem.Password = value
		}
	}
	return newItem
}

func NewCredentialStoreAdapter() tinystore.StoreItemAdapter {

	return tinystore.NewDefaultStoreItemAdapter(convert)
}

func NewCredentialStore() tinystore.Store {
	store := &tinystore.SimpleStore{Name: "CredentialStore"}
	tinystore.RegisterStoreAdapter(store, tinystore.NewDefaultStoreItemAdapter(convert))
	return store
}

