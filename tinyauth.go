package tinyauth

import (
	"net/http"
	"strings"
	"log"
	"github.com/D10221/tinyauth/config"
	"github.com/D10221/tinyauth/crypto"
	"github.com/D10221/tinyauth/tinystore"
	"github.com/D10221/tinyauth/encoder"
	"errors"
)
var (
	// ErrUnAuthorized Is Not Authorized
	ErrUnAuthorized = errors.New("unauthorized")
)
type Handler func(w http.ResponseWriter, r *http.Request);

type Authenticator func(u, p string) (bool, error);

type TinyAuth struct {
	Config          *config.TinyAuthConfig
	Criptico        crypto.Criptico
	CredentialStore tinystore.Store
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

func NewTinyAuth(config *config.TinyAuthConfig, astore tinystore.Store) *TinyAuth {

	if astore == nil {
		astore = &tinystore.SimpleStore{}
	}

	tAuth := &TinyAuth{
		Config:config,
		CredentialStore: astore,
		Criptico: &crypto.DefaultCriptico{config.Secret},
		Encoder: &encoder.DefaultEncoder{BasicScheme: config.BasicScheme},
	}
	return tAuth
}

func (t *TinyAuth) AuthFunc(username, password string) (bool, error) {

	found, e := t.CredentialStore.Find(tinystore.UserNameEquals(username))
	if e != nil && e != tinystore.ErrNotFound || !found.Valid() {
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

func (t *TinyAuth)NewCredential(username, password string) (*tinystore.Credential, error) {

	if t.Config.Secret == "" {
		return &tinystore.Credential{username, password}, nil
	}

	password, err := t.Criptico.Encrypt(password)
	if err != nil {
		return nil, err
	}

	return &tinystore.Credential{username, password}, nil
}

func (t *TinyAuth) GetRequestCredentials(r *http.Request) (*tinystore.Credential, error) {

	if t.Config.AuthorizationKey == "" {
		log.Println("Warning no AuthorizationKey")
	}

	auth := r.Header.Get(t.Config.AuthorizationKey)

	decoded, e := t.Encoder.Decode(auth)
	if e != nil {
		return &tinystore.Credential{}, e
	}
	parts := strings.Split(decoded, ":")
	if len(parts) < 2 {
		return &tinystore.Credential{}, e
	}
	return &tinystore.Credential{parts[0], parts[1]}, nil
}
var FormHasNoCredentials = errors.New("Form has no  credentials")

func (t *TinyAuth) GetFormCredentials(r *http.Request) (*tinystore.Credential, error) {

	e:= r.ParseForm()

	if e!=nil && e.Error()!= "missing form body"{
		return nil, e
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password =="" {
		username = r.PostForm.Get("username")
		password = r.PostForm.Get("password")
		// Still Empty ?
		if username == "" || password =="" {
			return nil,  FormHasNoCredentials
		}

	}

	credential := &tinystore.Credential{username, password}

	return credential, credential.Validate()
}

func (t TinyAuth) Authenticate(credential *tinystore.Credential) (bool, error) {

	if ! credential.Valid() {
		return false, tinystore.ErrInvalidCredential
	}

	ok, err := t.AuthFunc(credential.Username, credential.Password)

	if err == nil {
		return ok, nil
	}

	return false, err
}

func (t *TinyAuth) LoadConfig(path string) error {
	e := t.Config.LoadConfig(path)
	if e != nil {
		return e
	}
	if t.Config.Secret != "" {
		t.Criptico = &crypto.DefaultCriptico{t.Config.Secret}
	}
	if t.Config.BasicScheme != "" {
		t.Encoder = &encoder.DefaultEncoder{BasicScheme: t.Config.BasicScheme}
	}
	return nil
}

func (t *TinyAuth) EncryptPassword(in *tinystore.Credential) (*tinystore.Credential, error) {
	password, e := t.Criptico.Encrypt(in.Password)
	if e != nil {
		return nil, e
	}
	in.Password = password
	return in , nil
}

func (t *TinyAuth) Encode(credential *tinystore.Credential) (key string, value string) {
	return t.Config.AuthorizationKey, t.Encoder.Encode(credential.Username, credential.Password)
}
