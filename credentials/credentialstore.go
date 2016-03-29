package credentials

import (
	"github.com/D10221/tinystore"
)

// converts map[string]interface{} to Credential implementing StoreItem
var convert = func(item map[string]interface{}) tinystore.StoreItem {

	newItem := &Credential{}

	if keyValue, keyExists := item["Username"]; keyExists {
		if value, ok := keyValue.(string); ok {
			newItem.Username = value
		}
	}

	if keyValue, keyExists := item["Password"]; keyExists {
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

// Adapt Credential Filter to StoreItemFilter
func FilterAdapter(filter Filter) tinystore.Filter {

	return func(item tinystore.StoreItem) bool {
		c, e := ToCredential(item)
		if e != nil {
			return false
		}
		return filter(c)
	}
}

//  FindByKey shortcut
func FindByKey(store tinystore.Store, key string) (*Credential, error) {

	return ToCredentialMaybe(store.Find(CredentialNameFilter(key)))
}

// ToCredentialMaybe convert to Credential if (item,error) combo/tuple? is error free
func ToCredentialMaybe(item tinystore.StoreItem, e error) (*Credential, error) {

	if e != nil {
		return nil, e
	}
	return ToCredential(item)
}

// ToCredential Convert StoreItem to Credential
func ToCredential(item tinystore.StoreItem) (*Credential, error) {
	value, ok := item.(*Credential)
	if !ok {
		return nil, ErrInvalidCredential
	}
	return value, nil
}

// CredentialNameFilter return configured filter
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
func MutatorAdapter(m Mutator) tinystore.Mutator {
	return func(item tinystore.StoreItem) (tinystore.StoreItem, error) {
		value, ok := item.(*Credential)
		if ok {
			return m(value)
		}
		return value, ErrInvalidCredential
	}
}

// ForEach item in store satisfied by filter mutate item with mutator if no error found, if no filter provided , process all
func ForEach(store tinystore.Store, m Mutator, f Filter) error {
	if f == nil {
		f = Always
	}
	return  tinystore.ForEach(store, MutatorAdapter(m), FilterAdapter(f))
}

// All true if all items evaluate true and no error occurred on StoredItem conversion, first false or error stops loop
func All(store tinystore.Store, credentialFilter Filter) (bool, error) {
	for _, item := range store.All() {
		storeFilter := FilterAdapter(credentialFilter)
		if c,e := ToCredential(item); e!=nil || !storeFilter(c) {
			return false, e
		}
	}
	return true , nil
}

