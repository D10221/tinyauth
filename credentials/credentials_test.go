package credentials

import (
	"testing"
	"github.com/D10221/tinyauth/criptico"
)

func Test_Credentials(t *testing.T){

	store, ok := Credentials.(*SimpleCredentialStore)
	if !ok {
		t.Error("Credentials is Not a SimpleCredenetialStore")
	}
	store.Load(Credential{"admin", "P@55w0rd!"})

	cred := New("bob", "P@55w0rd!")

	encoded:= cred.Encode()

	decodedCredentials, err:= Decode(encoded)

	if err!= nil  {
		t.Error(err)
	}

	if decodedCredentials.Username != "bob" || decodedCredentials.Password!= "P@55w0rd!"{
		t.Error("It Doesn't Work")
	}

	if cred.Authenticate() {
		t.Error("shoudln't authorize")
	}

	if !AuthFunc("admin", "P@55w0rd!") {
		t.Error("Doesn't work")
	}

	if ! MustDecode(New("admin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It Doesn't work")
	}

	if  MustDecode(New("xxxadmin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It Doesn't work")
	}

	if  !ShouldDecode(New("admin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It should")
	}

	if  ShouldDecode(New("", "").Encode()).Authenticate(){
		t.Error("It shouldn't")
	}

	if  ShouldDecode((&Credential{}).Encode()).Authenticate(){
		t.Error("It shouldn't")
	}

	if  ShouldDecode(New("xxxadmin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It shoudn't")
	}

	if((&Credential{}).Valid()){
		t.Error("Shoudn't")
	}

	if(!(&Credential{"me", "1234"}).Valid()){
		t.Error("it Should")
	}
}

func Test_Encrypted_Credentials(t *testing.T){

	store, ok := Credentials.(*SimpleCredentialStore)
	if !ok {
		t.Error("Credentials is Not a *SimpleCredentialStore")
	}
	store.Load(Credential{"admin", "P@55w0rd!"})

	secret := []byte("ABRACADABRA12345")

	cred := New("bob", criptico.Encrypt(secret, "P@55w0rd!"))

	encoded:= cred.Encode()

	decodedCredentials, err:= Decode(encoded)

	if err!= nil  {
		t.Error(err)
	}

	if decodedCredentials.Username != "bob" || criptico.Decrypt(secret, decodedCredentials.Password)!= "P@55w0rd!"{
		t.Errorf("It Doesn't Work, %s", decodedCredentials.Password)
	}

	if cred.Authenticate() {
		t.Error("shoudln't authorize")
	}

	if !AuthFunc("admin", "P@55w0rd!") {
		t.Error("Doesn't work")
	}

	if ! MustDecode(New("admin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It Doesn't work")
	}

	if  MustDecode(New("xxxadmin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It Doesn't work")
	}

	if  !ShouldDecode(New("admin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It should")
	}

	if  ShouldDecode(New("", "").Encode()).Authenticate(){
		t.Error("It shouldn't")
	}

	if  ShouldDecode((&Credential{}).Encode()).Authenticate(){
		t.Error("It shouldn't")
	}

	if  ShouldDecode(New("xxxadmin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It shoudn't")
	}

	if((&Credential{}).Valid()){
		t.Error("Shoudn't")
	}

	if(!(&Credential{"me", "1234"}).Valid()){
		t.Error("it Should")
	}
}

type _TestSTore struct {

}

func (store *_TestSTore) All() []Credential{
	return []Credential{ Credential{"admin", "P@55w0rd!"} }
}

func Test_Replace_CredentialStore(t *testing.T){
	Credentials = &_TestSTore{}
	if Credentials.All()[0].Username != "admin" {
		t.Error("Doen't work")
	}
}
