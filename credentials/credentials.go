package credentials

import (

	"github.com/D10221/tinyauth/criptico"
	"github.com/D10221/tinyauth/config"
)


func MustDecode(auth string) *Credential {
	credentials, err := Decode(auth)
	if err != nil {
		panic(err)
	}
	return credentials
}

func ShouldDecode(auth string) *Credential {
	credentials, _ := Decode(auth)
	return credentials
}

type Authenticator func(u, p string) bool;

var AuthFunc Authenticator = nil

var Credentials CredentialStore = &SimpleCredentialStore{};

type CredentialStore interface {
	All() []Credential
}

func init() {
	AuthFunc = func(username, password string) bool {
		for _, credential := range Credentials.All()[:] {
			// ternary operator
			currentPassword := func() string{
				if config.Current.Secret == "" {
					return credential.Password
				}
				return criptico.Decrypt([]byte(config.Current.Secret), credential.Password)
			}()

			if credential.Username == username &&
			currentPassword == password {
				return true
			}
		}
		return false
	};
}