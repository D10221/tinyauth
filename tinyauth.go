package tinyauth

import (
	"net/http"
	"strings"
	"log"
	"github.com/D10221/tinyauth/config"
	"github.com/D10221/tinyauth/crypto"
	"github.com/D10221/tinystore"
	"github.com/D10221/tinyauth/encoder"
	"errors"
	"github.com/D10221/tinyauth/credentials"
)
var (
	// ErrUnAuthorized Is Not Authorized
	ErrUnAuthorized = errors.New("unauthorized")
)

// Handler
type Handler func(w http.ResponseWriter, r *http.Request);


// Authenticator
type Authenticator func(u, p string) (bool, error);

// TinyAuth glue...
type TinyAuth struct {
	Config          *config.TinyAuthConfig
	Criptico        crypto.Criptico
	CredentialStore tinystore.Store
	Encoder         encoder.Encoder
}

// RequireAuthentication main purpose of this whole exercise ,will end up as middleware/plugin
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

// NewTinyAuth default TinyAuth fty, optional store, can be nil , will use default implementation
func NewTinyAuth(config *config.TinyAuthConfig, credentialStore tinystore.Store) *TinyAuth {

	if credentialStore == nil {
		credentialStore = credentials.NewCredentialStore()
		tinystore.RegisterStoreAdapter(credentialStore, credentials.NewCredentialStoreAdapter())
	}

	tAuth := &TinyAuth{
		Config:config,
		CredentialStore: credentialStore,
		Criptico: &crypto.DefaultCriptico{config.Secret},
		Encoder: &encoder.DefaultEncoder{BasicScheme: config.BasicScheme},
	}
	return tAuth
}

// AuthFunc
func (t *TinyAuth) AuthFunc(username, password string) (bool, error) {
	found, e:= credentials.FindByKey(t.CredentialStore, username)
	if e!= nil {
		return false, e
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

// NewCredential
func (t *TinyAuth)NewCredential(username, password string) (*credentials.Credential, error) {

	if t.Config.Secret == "" {
		return &credentials.Credential{username, password}, nil
	}

	password, err := t.Criptico.Encrypt(password)
	if err != nil {
		return nil, err
	}

	return &credentials.Credential{username, password}, nil
}

// GetRequestCredentials
func (t *TinyAuth) GetRequestCredentials(r *http.Request) (*credentials.Credential, error) {

	if t.Config.AuthorizationKey == "" {
		log.Println("Warning no AuthorizationKey")
	}

	auth := r.Header.Get(t.Config.AuthorizationKey)

	decoded, e := t.Encoder.Decode(auth)
	if e != nil {
		return &credentials.Credential{}, e
	}
	parts := strings.Split(decoded, ":")
	if len(parts) < 2 {
		return &credentials.Credential{}, e
	}
	return &credentials.Credential{parts[0], parts[1]}, nil
}

// FormHasNoCredentials
var ErrFormHasNoCredentials = errors.New("Form has no  credentials")

// GetFormCredentials
func (t *TinyAuth) GetFormCredentials(r *http.Request) (*credentials.Credential, error) {

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
			return nil, ErrFormHasNoCredentials
		}

	}

	credential := &credentials.Credential{username, password}

	return credential, credential.Validate()
}

// Authenticate
func (t TinyAuth) Authenticate(credential *credentials.Credential) (bool, error) {

	if ! credential.Valid() {
		return false, tinystore.ErrInvalidStoreItem
	}

	ok, err := t.AuthFunc(credential.Username, credential.Password)

	if err == nil {
		return ok, nil
	}

	return false, err
}

// LoadConfig
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

// EncryptPassword
func (t *TinyAuth) EncryptPassword(in *credentials.Credential) (*credentials.Credential, error) {

	password, e := t.Criptico.Encrypt(in.Password)
	if e != nil {
		return nil, e
	}
	in.Password = password
	return in , nil
}

// Encode
func (t *TinyAuth) Encode(credential *credentials.Credential) (key string, value string) {
	return t.Config.AuthorizationKey, t.Encoder.Encode(credential.Username, credential.Password)
}
