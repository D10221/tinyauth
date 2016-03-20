package credentials

import (
	"strings"
	"log"
	"github.com/D10221/tinyauth/encoding"
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

func (c *Credential) Authenticate() bool {
	if AuthFunc == nil {
		log.Fatal("No Authenticator found")
		return false
	}
	return AuthFunc(c.Username, c.Password)
}
