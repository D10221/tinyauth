package store

import (
	"io/ioutil"
	"encoding/json"
	"errors"
)

type CredentialStore interface {
	All() []*Credential
	Load(credentials ...*Credential)
	FindUser(userName string) *Credential
	LoadJson(path string) error
	Add(credential *Credential) error
	Remove(credential *Credential) error
	// Update all with Func
	Update(transform func(in *Credential) *Credential) error
}

type CredentialStoreError struct {
	message string
	Code int
}

func (e *CredentialStoreError) Error() string {
	return e.message;
}

func NewCredentialStoreError(message string , code int ) *CredentialStoreError {
	return &CredentialStoreError{message, code}
}

type SimpleCredentialStore struct {
	all []*Credential
}

func (store *SimpleCredentialStore) All() []*Credential {
	return store.all
}

func (store *SimpleCredentialStore) FindUser(userName string) *Credential {
	found := &Credential{}
	for _, credential := range store.All()[:] {
		if userName == credential.Username {
			found = credential
			break
		}
	}
	return found
}


func (store *SimpleCredentialStore) FindValidUser(userName string) (*Credential, error) {
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
	return found, NewCredentialStoreError("Credential Not Found", 404)
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

var InvalidCredential = errors.New("Invalid Credential")
var NotFound = errors.New("Credential not found")
var AlreadyExists = errors.New("Credential Already Exists")

func (store *SimpleCredentialStore) Add(credential *Credential) error {
	if !credential.Valid() {
		return InvalidCredential
	}
	found := store.FindUser(credential.Username)
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
	result:= make([]*Credential,0)
	found:= false
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

func (store *SimpleCredentialStore) Update(transform func(in *Credential) *Credential) error {
	all := make([]*Credential,0)
	for _, item := range store.all {
		result:= transform(item)
		if result.Valid(){
			all = append(all, result)
		}
	}
	store.all = all
	return nil
}

