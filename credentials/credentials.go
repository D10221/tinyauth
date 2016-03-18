package credentials

import (
	"github.com/D10221/tinyauth/encoding"
	"strings"
)

type Credentials struct {
	Username string
	Password string
}

func New(username, password string) *Credentials {
	return &Credentials{Username:username, Password: password}
}

func (credentials *Credentials) Encode() string {
	return encoding.EncodeWithSchema(credentials.Username, credentials.Password)
}

func Decode(auth string) (*Credentials, error) {
	c := &Credentials{}
	decoded, e, ok := encoding.Decode(auth)
	if ok {
		parts := strings.Split(decoded, ":")
		if len(parts) < 2 {
			// t.Error("Bad Decoding")
			return c, e
		}
		c.Username = parts[0]
		c.Password = parts[1]
	}
	return c, nil
}

func AuthFunc(u, p string) bool {
	return u == "u" && p == "p"
}

func (c *Credentials)Authenticate() bool {
	return AuthFunc(c.Username, c.Password)
}