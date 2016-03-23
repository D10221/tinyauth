package tinyauth

import (
	"net/http"
	"strings"
	"errors"
	"log"
	"github.com/D10221/tinyauth/config"
	"github.com/D10221/tinyauth/crypto"
	"github.com/D10221/tinyauth/store"
	"github.com/D10221/tinyauth/encoder"
)

type Handler func(w http.ResponseWriter, r *http.Request);

type Authenticator func(u, p string) (bool, error);

type TinyAuth struct {
	Config          *config.TinyAuthConfig
	Criptico        crypto.Criptico
	CredentialStore store.CredentialStore
	Encoder         encoder.Encoder
}

func (t TinyAuth) RequireAuthentication(handler Handler) Handler {
	// ...
	return func(w http.ResponseWriter, r *http.Request) {

		credential, err := t.GetRequestCredentials(r)

		if err != nil {
			log.Printf("RequireAuthentication: Error: %s \n", err.Error())
		}

		ok, err := t.Authenticate(credential)

		if err != nil {
			log.Printf("RequireAuthentication: Error: %s", err.Error())
			http.Error(w, "unauthorized", 401)
			return
		}

		if !ok {
			http.Error(w, "unauthorized", 401)
			return

		}

		handler(w, r)
	}
}

func NewTinyAuth(config *config.TinyAuthConfig) *TinyAuth {

	tAuth := &TinyAuth{
		Config:config,
		CredentialStore: &store.SimpleCredentialStore{},
		Criptico: &crypto.DefaultCriptico{config.Secret},
		Encoder: &encoder.DefaultEncoder{BasicScheme: config.BasicScheme},
	}
	return tAuth
}


func (t *TinyAuth) AuthFunc(username, password string) (bool, error) {

	found, e := t.CredentialStore.FindByUserName(username)
	if e!=nil && e!= store.NotFound  || !found.Valid() {
		return false, nil
	}
	if t.Config.Secret == "" {
		return found.Username == username && found.Password == password, nil
	}

	currentPassword, err := t.Criptico.Decrypt(found.Password)
	if err != nil {
		return false, err
	}
	return found.Username == username && currentPassword == password, nil
}

func (t *TinyAuth)NewCredential(username, password string) (*store.Credential, error) {

	if t.Config.Secret == "" {
		return &store.Credential{username, password}, nil
	}

	password, err := t.Criptico.Encrypt(password)
	if err != nil {
		return nil, err
	}

	return &store.Credential{username, password} , nil
}

func (t *TinyAuth) GetRequestCredentials(r *http.Request) (*store.Credential, error) {

	if t.Config.AuthorizationKey == "" {
		log.Println("Warning no AuthorizationKey")
	}

	auth := r.Header.Get(t.Config.AuthorizationKey)

	decoded, e := t.Encoder.Decode(auth)
	if e != nil {
		return &store.Credential{}, e
	}
	parts := strings.Split(decoded, ":")
	if len(parts) < 2 {
		return &store.Credential{}, e
	}
	return &store.Credential{parts[0], parts[1]}, nil
}

func (t *TinyAuth) GetFormCredentials(r *http.Request) (*store.Credential, error) {

	if t.Config.AuthorizationKey == "" {
		log.Println("Warning no AuthorizationKey")
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	return &store.Credential{username, password}, nil
}

func (t TinyAuth) Authenticate(credential *store.Credential) (bool, error) {

	if ! credential.Valid() {
		return false, errors.New("Credential not valid")
	}

	ok, err := t.AuthFunc(credential.Username, credential.Password)

	if err == nil { return ok, nil 	}

	return false, err
}



func (t *TinyAuth) LoadConfig(path string) error{
	e:= t.Config.LoadConfig(path)
	if e!= nil {
		return e
	}
	if t.Config.Secret!= "" {
		t.Criptico = &crypto.DefaultCriptico{t.Config.Secret}
	}
	if t.Config.BasicScheme != "" {
		t.Encoder = &encoder.DefaultEncoder{BasicScheme: t.Config.BasicScheme}
	}
	return nil
}

