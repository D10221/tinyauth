package store

type Credential struct {
	Username string
	Password string
}


// Must contain values
func (cred *Credential) Valid() bool {
	return cred.Username != "" && cred.Password != ""
}

type CredentialComparison func (a *Credential,other *Credential) bool ;

type CredentialValidation func (a *Credential) bool ;

func AreNamesEqual(a *Credential,b *Credential) bool {
	return a.Username == b.Username
}

func AreEqual (a *Credential,b *Credential) bool {
	if a ==nil || b == nil {
		return a == b
	}
	return a.Username == b.Username && a.Password ==b.Password
}

func ByName(name string) CredentialValidation {
	validator:= func (c *Credential) bool {
		return c!=nil && c.Username == name
	}
	return validator
}
