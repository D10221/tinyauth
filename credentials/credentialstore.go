package credentials

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"os"
)

type SimpleCredentialStore struct {
	all [] Credential
}

func(store *SimpleCredentialStore) All() []Credential {
	return store.all
}

// "testdata/credentials.json"
func (store *SimpleCredentialStore) LoadJsonFromRelativePath(path string) {

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

func (store *SimpleCredentialStore) Load(credentials ...Credential){
	store.all = make([]Credential, 0)
	for _,credential := range credentials[:]{
		store.all = append(store.all, credential)
	}
}


