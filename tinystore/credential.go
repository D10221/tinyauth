package tinystore

type Credential struct {
	Username string
	Password string
}


// Must contain valid values
func (cred *Credential) Valid() bool {
	return cred != nil &&  cred.Username != "" && cred.Password != ""
}

// validate credential values, return error if Not
func (cred *Credential) Validate() error {
	if cred!=nil &&  cred.Valid() {
		return nil
	}
	return ErrInvalidCredential
}

func (cred *Credential) Equals(other *Credential) bool {
	if cred == nil || other == nil {
		return other == cred
	}
	return cred.Username == other.Username
}

type CredentialComparison func (a *Credential,other *Credential) bool ;

type Filter func (a *Credential) bool ;

func NotFilter(filter Filter) Filter {
	return func(c *Credential) bool {
		return !filter(c)
	}
}

type Mutator  func(in *Credential) (*Credential, error );

var Always Filter = func( a *Credential) bool {return  true }

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

func UserNameEquals(name string) Filter {
	validator:= func (c *Credential) bool {
		return c!=nil && c.Username == name
	}
	return validator
}
