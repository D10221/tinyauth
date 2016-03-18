package credentials

import (
	"github.com/D10221/tinyauth/encoding"
	"strings"
	"log"
)

type Credentials struct {
	Username string
	Password string
}

func New(username, password string) *Credentials {
	return &Credentials{Username:username, Password: password}
}

// Must contain values
func (cred *Credentials) Valid() bool {
	return cred.Username != "" && cred.Password != ""
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

func MustDecode(auth string) *Credentials {
	credentials, err := Decode(auth)
	if err != nil {
		panic(err)
	}
	return credentials
}

func ShouldDecode(auth string) *Credentials {
	credentials, _ := Decode(auth)
	return credentials
}

var AllCredentials  []Credentials

type Authenticator func (u, p string) bool;

var authenticator Authenticator = nil ;

func (c *Credentials) Authenticate() bool {
	if authenticator == nil {
		log.Fatal("No Authenticator found")
		return false
	}
	return  authenticator(c.Username, c.Password)
}

func init(){
	authenticator = func (u, p string) bool {
		for _, credential := range AllCredentials[:]{
			if credential.Username == u && credential.Password == p {
				return true
			}
		}
		return false
	};
}