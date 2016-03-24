package store

import (
	"testing"
)


func Test_Credentials(t *testing.T) {

	cred := &Credential{"bob", "P@55w0rd!"}
	if cred == nil {
		t.Error("WTF")
	}
}

func Test_Is_Valid_Credential(t *testing.T){

	if ((&Credential{}).Valid()) {
		t.Error("Shoudn't")
	}

	if (!(&Credential{"me", "1234"}).Valid()) {
		t.Error("it Should")
	}
}

func Test_AreEqual(t *testing.T) {

	if AreEqual(&Credential{"A", "x"},&Credential{"A", "x1"}) {
		t.Error("They Are Not Equal")
	}
	if AreEqual(&Credential{"B", "x"},&Credential{"A", "x"}) {
		t.Error("They Are Not Equal")
	}
	if AreEqual(&Credential{"A", "x"},&Credential{}) {
		t.Error("They Are Not Equal")
	}
	if !AreEqual(&Credential{"A", "x"},&Credential{"A", "x"}) {
		t.Error("They Are Equal")
	}
	if !AreEqual(&Credential{},&Credential{}) {
		t.Error("They Are Equal")
	}
	if AreEqual(nil,&Credential{}) {
		t.Error("They Are Not Equal")
	}
	if AreEqual(&Credential{},nil) {
		t.Error("They Are Not Equal")
	}
	if !AreEqual(nil,nil) {
		t.Error("They Are Equal")
	}
}

func Test_ByName(t *testing.T){
	ok:= AreNamesEqual(&Credential{}, &Credential{})
	if !ok {
		t.Error("It Should")
	}
	ok= AreNamesEqual(&Credential{}, &Credential{"A", "X"})
	if ok {
		t.Error("It Shouldn't")
	}
	ok= AreNamesEqual(nil, &Credential{"A", "X"})
	if ok {
		t.Error("It Shouldn't")
	}
	ok= AreNamesEqual(&Credential{"A", "X"}, nil)
	if ok {
		t.Error("It Shouldn't")
	}
}



