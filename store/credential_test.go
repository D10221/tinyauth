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



