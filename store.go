package tinyauth

import (
	"io/ioutil"
	"errors"
	"encoding/json"
)

type CredentialStore interface {
	All() []Credential
	Load(credentials ...Credential)
	FindUser(userName string) (Credential, error)
	LoadJson(path string) error
}


type SimpleCredentialStore struct {
	all []Credential
}

func (store *SimpleCredentialStore) All() []Credential {
	return store.all
}

func (store *SimpleCredentialStore) FindUser(userName string) (Credential, error) {
	var found Credential
	for _, credential := range store.All()[:] {
		if userName == credential.Username {
			found = credential
			break
		}
	}
	if found.Valid() {
		return found, nil
	}
	return found, errors.New("Credential Not Found")
}

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

func (store *SimpleCredentialStore) Load(credentials ...Credential) {
	store.all = make([]Credential, 0)
	for _, credential := range credentials[:] {
		store.all = append(store.all, credential)
	}
}

