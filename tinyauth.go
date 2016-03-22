package tinyauth

import (
	"net/http"
	"strings"
	"errors"
	"log"
)

type Handler func(w http.ResponseWriter, r *http.Request);

type Authenticator func(u, p string) (bool, error);

type TinyAuth struct {
	Config          *TinyAuthConfig
	Criptico        Criptico
	CredentialStore CredentialStore
	Encoder         Encoder
}

func (t TinyAuth) RequireAuthentication(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {

		credential, err := t.GetCredential(r)

		if err != nil {
			log.Printf("Error: %s \n", err.Error())
		}

		if ok, err := t.AuthFunc(credential.Username, credential.Password); ok && (err == nil ) {
			handler(w, r)
		} else {
			if err != nil {
				//log.Printf("Error: %s", err.Error())
			}
			http.Error(w, "unauthorized", 401)
			return
		}
	}
}

func NewTinyAuth(config *TinyAuthConfig) *TinyAuth {

	tAuth := &TinyAuth{
		Config:config,
		CredentialStore: &SimpleCredentialStore{},
		Criptico: &DefaultCriptico{config.Secret},
		Encoder: &DefaultEncoder{config: config},
	}

	return tAuth
}

func (t *TinyAuth) AuthFunc(username, password string) (bool, error )  {
	found, err := t.CredentialStore.FindUser(username)
	if err != nil {
	return false, err
	}
	if t.Config.Secret == "" {
	return found.Username == username && found.Password == password, nil
	}
	currentPassword, err := t.Criptico.Decrypt(found.Password)
	if err!= nil {
	return false, err
	}
	return found.Username == username && currentPassword == password, nil
}

func (t *TinyAuth)NewCredential(username, password string) *Credential {

	if t.Config.Secret == "" {
		return &Credential{username, password}
	}

	password, err := t.Criptico.Encrypt(password)
	if err != nil {
		panic(err)
	}

	return &Credential{username, password}
}

func (t *TinyAuth) GetCredential(r *http.Request) (*Credential, error) {

	if t.Config.AuthorizationKey == "" {
		//log.Println("Warning no config.Current.AuthorizationKey")
	}

	auth := r.Header.Get(t.Config.AuthorizationKey)

	decoded, e := t.Encoder.Decode(auth)
	if e != nil {
		return &Credential{}, e
	}
	parts := strings.Split(decoded, ":")
	if len(parts) < 2 {
		return &Credential{}, e
	}
	return &Credential{parts[0], parts[1]}, nil
}


func (t TinyAuth) Authenticate(c *Credential) (bool, error) {
	if t.AuthFunc == nil {
		return false, errors.New(("No Authenticator found"))
	}
	return t.AuthFunc(c.Username, c.Password)
}

