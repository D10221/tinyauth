package credentials

import (
	"strings"
	"log"
	"github.com/D10221/tinyauth/encoding"
	"github.com/D10221/tinyauth/criptico"
	"github.com/D10221/tinyauth/config"
)

type Credential struct {
	Username string
	Password string
}

func New(username, password string) *Credential {
	return &Credential{Username:username, Password: password}
}

// Must contain values
func (cred *Credential) Valid() bool {
	return cred.Username != "" && cred.Password != ""
}

func (credentials *Credential) Encode() string {
	return encoding.EncodeWithSchema(credentials.Username, credentials.Password)
}

func Decode(auth string) (*Credential, error) {
	c := &Credential{}
	decoded, e := encoding.Decode(auth)
	if e == nil {
		parts := strings.Split(decoded, ":")
		if len(parts) < 2 {
			// t.Error("Bad Decoding")
			return c, e
		}
		c.Username = parts[0]
		c.Password = parts[1]
	}
	return c, e
}

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



func (c *Credential) Authenticate() bool {
	if AuthFunc == nil {
		log.Fatal("No Authenticator found")
		return false
	}
	return AuthFunc(c.Username, c.Password)
}

type Authenticator func(u, p string) bool;

var AuthFunc Authenticator = nil;

var Credentials = &CredentialStore{}

func init() {
	AuthFunc = func(username, password string) bool {
		for _, credential := range Credentials.All()[:] {
			currentPassword := func() string{
				// TODO Refactor, circular reference ?
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