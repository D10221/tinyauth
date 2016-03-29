package credentials_test

import (
	"testing"
	"github.com/D10221/tinyauth/credentials"
)

func Test_Credentials(t *testing.T) {

	cred := &credentials.Credential{"bob", "P@55w0rd!"}
	if cred == nil {
		t.Error("WTF")
	}
}

func Test_Is_Valid_Credential(t *testing.T){

	if ((&credentials.Credential{}).Valid()) {
		t.Error("Shoudn't")
	}

	if (!(&credentials.Credential{"me", "1234"}).Valid()) {
		t.Error("it Should")
	}
}

func Test_AreEqual(t *testing.T) {

	if credentials.AreEqual(&credentials.Credential{"A", "x"},&credentials.Credential{"A", "x1"}) {
		t.Error("They Are Not Equal")
	}
	if credentials.AreEqual(&credentials.Credential{"B", "x"},&credentials.Credential{"A", "x"}) {
		t.Error("They Are Not Equal")
	}
	if credentials.AreEqual(&credentials.Credential{"A", "x"},&credentials.Credential{}) {
		t.Error("They Are Not Equal")
	}
	if !credentials.AreEqual(&credentials.Credential{"A", "x"},&credentials.Credential{"A", "x"}) {
		t.Error("They Are Equal")
	}
	if !credentials.AreEqual(&credentials.Credential{},&credentials.Credential{}) {
		t.Error("They Are Equal")
	}
	if credentials.AreEqual(nil,&credentials.Credential{}) {
		t.Error("They Are Not Equal")
	}
	if credentials.AreEqual(&credentials.Credential{},nil) {
		t.Error("They Are Not Equal")
	}
	if !credentials.AreEqual(nil,nil) {
		t.Error("They Are Equal")
	}
}

func Test_ByName(t *testing.T){
	ok:= credentials.AreNamesEqual(&credentials.Credential{}, &credentials.Credential{})
	if !ok {
		t.Error("It Should")
	}
	ok= credentials.AreNamesEqual(&credentials.Credential{}, &credentials.Credential{"A", "X"})
	if ok {
		t.Error("It Shouldn't")
	}
	ok= credentials.AreNamesEqual(nil, &credentials.Credential{"A", "X"})
	if ok {
		t.Error("It Shouldn't")
	}
	ok= credentials.AreNamesEqual(&credentials.Credential{"A", "X"}, nil)
	if ok {
		t.Error("It Shouldn't")
	}
}



