package crypto

import (
	"testing"
)

func Test_Criptico(t *testing.T){
	var crypt Criptico = &DefaultCriptico{Key: "0123456789ABCDEF"}
	expected := "hello"
	encrypted, e:= crypt.Encrypt(expected)
	if e != nil {
		t.Error(e)
	}
	decrypted, e := crypt.Decrypt(encrypted)
	if e != nil {
		t.Error(e)
	}
	if expected!= decrypted {
		t.Fail()
	}
}

func Test_Criptico_Bad_Key(t *testing.T){

	var crypt Criptico = &DefaultCriptico{Key: "0"}

	expected := "hello"
	encrypted, e:= crypt.Encrypt(expected)
	if e == nil {
		t.Error("It Should fail")
	}
	decrypted, e := crypt.Decrypt(encrypted)
	if e == nil {
		t.Error("it should Fail")
	}
	if encrypted != "" {
		t.Error("Should be empty")
	}
	if decrypted != "" {
		t.Error("Should be empty")
	}
}

