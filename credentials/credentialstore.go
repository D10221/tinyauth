package credentials

import (
	"github.com/D10221/tinystore"
)

var convert = func(item map[string]interface{}) tinystore.StoreItem {

	newItem:= &Credential{}

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

// NewCredentialStoreAdapter local Credential adapter
func NewCredentialStoreAdapter() tinystore.StoreItemAdapter {

	return tinystore.NewDefaultStoreItemAdapter(convert)
}

// NewCredentialStore default Fty
func NewCredentialStore() tinystore.Store {
	store := &tinystore.SimpleStore{Name: "CredentialStore"}
	tinystore.RegisterStoreAdapter(store, tinystore.NewDefaultStoreItemAdapter(convert))
	return store
}

// Find by filter: credentials.Filter
func Find(store tinystore.Store, filter Filter) (*Credential, error) {
	return ToCredentialMaybe(store.Find(FilterAdapter(filter)))
}

func FilterAdapter (filter Filter) tinystore.Filter {

	return func(item tinystore.StoreItem)bool{
		c, e:= ToCredential(item )
		if e!=nil {
			return false
		}
		return filter(c)
	}
}

func FindByKey(store tinystore.Store, key string) (*Credential, error) {

	return ToCredentialMaybe(store.Find(CredentialNameFilter(key)))
}


func ToCredentialMaybe(item tinystore.StoreItem,e error) (*Credential, error) {

	if e != nil {
		return nil , e
	}
	return ToCredential(item)
}

func ToCredential(item tinystore.StoreItem) (*Credential, error) {
	value, ok := item.(*Credential)
	if !ok {
		return nil, ErrInvalidCredential
	}
	return value, nil
}

// CredentialNameFilter
func CredentialNameFilter(name string) tinystore.Filter {
	return func(item tinystore.StoreItem) bool {

		value, ok := item.(*Credential)
		if !ok {
			return false
			//panic("Not a credential")
		}
		return value.Username == name
	}
}

// ForEach(app.Auth.CredentialStore, app.Auth.EncryptPassword , credentials.Always)
func MutatorAdapter (m Mutator) tinystore.Mutator {
	return func(item tinystore.StoreItem) (tinystore.StoreItem, error) {
		value, ok := item.(*Credential)
		if ok  {
			return m(value)
		}
		return value, ErrInvalidCredential
	}
}

func ForEach(store tinystore.Store,m Mutator, f Filter) {
	if f == nil {
		f = Always
	}
	tinystore.ForEach(store, MutatorAdapter(m), FilterAdapter(f))
}

