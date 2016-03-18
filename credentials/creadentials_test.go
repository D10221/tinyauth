package credentials

import (
	"testing"
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
	"path/filepath"
)

func setupCredentials(){

	dir, e:= os.Getwd()
	if e != nil{
		log.Fatal(e)
	}
	bytes, e := ioutil.ReadFile(filepath.Join(dir, "../testdata/credentials.json"))
	if e!= nil {
		panic(e)
	}
	e = json.Unmarshal(bytes, &AllCredentials)
	if e!= nil{
		panic(e)
	}
}

func Test_Credentials(t *testing.T){

	setupCredentials()

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

	if !authenticator("admin", "P@55w0rd!") {
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

	if  ShouldDecode((&Credentials{}).Encode()).Authenticate(){
		t.Error("It shouldn't")
	}

	if  ShouldDecode(New("xxxadmin", "P@55w0rd!").Encode()).Authenticate(){
		t.Error("It shoudn't")
	}

	if((&Credentials{}).Valid()){
		t.Error("Shoudn't")
	}

	if(!(&Credentials{"me", "1234"}).Valid()){
		t.Error("it Should")
	}
}
