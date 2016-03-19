package credentials

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"os"
)

type CredentialStore struct {
	all []Credential
}

func(store *CredentialStore) All() []Credential {
	return store.all
}

// "testdata/credentials.json"
func (store *CredentialStore) LoadJsonFromRelativePath(path string) {

	dir, e := os.Getwd()
	if e != nil {
		// no credentials no app
		panic(e)
	}
	bytes, e := ioutil.ReadFile(filepath.Join(dir, path))
	if e != nil {
		// no credentials no app
		panic(e)
	}
	e = json.Unmarshal(bytes, &store.all)
	if e != nil {
		// no credentials no app
		panic(e)
	}
}

func (store *CredentialStore) Load(credentials ...Credential){
	store.all = make([]Credential, 0)
	for _,credential := range credentials[:]{
		store.all = append(store.all, credential)
	}
}

