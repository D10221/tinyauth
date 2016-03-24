package store

type Credential struct {
	Username string
	Password string
}


// Must contain valid values
func (cred *Credential) Valid() bool {
	return cred.Username != "" && cred.Password != ""
}

// validate credential values, return error if Not
func (cred *Credential) Validate() error {
	if cred!=nil &&  cred.Valid() {
		return nil
	}
	return InvalidCredential
}

type CredentialComparison func (a *Credential,other *Credential) bool ;

type CredentialFilter func (a *Credential) bool ;

type CredentialMutator  func(in *Credential) (*Credential, error );

var Always CredentialFilter = func( a *Credential) bool {return  true }

func AreNamesEqual(a *Credential,b *Credential) bool {
	if a ==nil || b == nil { return false }
	return a.Username == b.Username
}

func AreEqual (a *Credential,b *Credential) bool {
	if a ==nil || b == nil {
		return a == b
	}
	return a.Username == b.Username && a.Password ==b.Password
}

func ByName(name string) CredentialFilter {
	validator:= func (c *Credential) bool {
		return c!=nil && c.Username == name
	}
	return validator
}
