package store

import (
	"io/ioutil"
	"encoding/json"
	"errors"
)

type CredentialStore interface {
	All() []*Credential
	Load(credentials ...*Credential)
	FindByUserName(userName string) (*Credential, error)
	FindBy(c CredentialFilter) (*Credential, error)
	LoadJson(path string) error
	Add(credential *Credential) error
	Remove(credential *Credential) error
	RemoveWhere(f CredentialFilter) error
	// Update all with Func
	UpdateAll(transform func(in *Credential) *Credential) error
	UpdateWhere(f CredentialFilter,transform CredentialMutator) error
}

type CredentialStoreError struct {
	message string
	Code    int
}

type SimpleCredentialStore struct {
	all []*Credential
}

func (store *SimpleCredentialStore) All() []*Credential {
	return store.all
}

/*func (store *SimpleCredentialStore) FindUser(userName string) *Credential {
	found := &Credential{}
	for _, credential := range store.All()[:] {
		if userName == credential.Username {
			found = credential
			break
		}
	}
	return found
}*/


func (store *SimpleCredentialStore) FindByUserName(userName string) (*Credential, error) {
	found := &Credential{}
	for _, credential := range store.All()[:] {
		if userName == credential.Username {
			found = credential
			break
		}
	}
	if found.Valid() {
		return found, nil
	}

	return found, NotFound
}
/*
func FilterCredentialStoreErrorNotFound(err error) error {
	if err == nil { return nil }
	if e, isCredentialStoreError := err.(*CredentialStoreError); isCredentialStoreError && e.Code == 404 {
		return nil
	}
	return err
}*/
func (store *SimpleCredentialStore) LoadJson(path string) error {

	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	e = json.Unmarshal(bytes, &store.all)
	if e != nil {
		return e
	}
	return nil
}

func (store *SimpleCredentialStore) Load(credentials ...*Credential) {
	store.all = make([]*Credential, 0)
	for _, credential := range credentials[:] {
		store.all = append(store.all, credential)
	}
}

var (
	InvalidCredential = errors.New("Invalid Credential")
	NotFound = errors.New("Credential not found")
	AlreadyExists = errors.New("Credential Already Exists")
)

func (store *SimpleCredentialStore) Add(credential *Credential) error {
	if !credential.Valid() {
		return InvalidCredential
	}
	found, e := store.FindByUserName(credential.Username)
	if e != nil && e != NotFound {
		return e
	}
	if found.Valid() {
		return AlreadyExists
	}
	store.all = append(store.all, credential)
	return nil
}

func (store *SimpleCredentialStore) Remove(credential *Credential) error {
	if !credential.Valid() {
		return InvalidCredential
	}
	result := make([]*Credential, 0)
	found := false
	for _, item := range store.all {
		if item.Username == credential.Username {
			found = true
			continue
		}
		result = append(result, item)
	}
	if !found {
		return NotFound
	}
	store.all = result
	return nil
}

func (store *SimpleCredentialStore) RemoveWhere(f CredentialFilter) error {

	result := make([]*Credential, 0)

	found := false
	for _, item := range store.all {
		if f(item){
			found = true
			continue
		}
		result = append(result, item)
	}
	if !found {
		return NotFound
	}
	store.all = result
	//TODO:
	return nil
}

func (store *SimpleCredentialStore) UpdateAll(transform func(in *Credential) *Credential) error {
	all := make([]*Credential, 0)
	var e error = nil
	for _, item := range store.all {
		result := transform(item)
		if result.Valid() {
			all = append(all, result)
		} else {
			e = InvalidCredential
		}
	}
	if e == nil {
		store.all = all
	}

	return e
}

func (store *SimpleCredentialStore) UpdateWhere(filter CredentialFilter,transform CredentialMutator) error {

	all := make([]*Credential, 0)
	var e error = NotFound
	for _, item := range store.all {
		if filter(item) {
			result, err := transform(item)
			all = append(all, result)
			e = err
		}
	}
	if e == nil {
		store.all = all
	}
	return e
}

func (store *SimpleCredentialStore) FindBy(c CredentialFilter) (*Credential, error) {
	found := &Credential{}
	for _, item := range store.all {
		if c(item) {
			found = item
			break
		}
	}
	if found.Valid() {
		return found, nil
	}
	return found, NotFound
}

