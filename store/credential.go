package store

type Credential struct {
	Username string
	Password string
}


// Must contain values
func (cred *Credential) Valid() bool {
	return cred.Username != "" && cred.Password != ""
}
