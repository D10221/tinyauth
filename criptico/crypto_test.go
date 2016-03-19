package criptico

import (
	"testing"
)

func Test_Crypto(t *testing.T){

	originalText := "encrypt this"

	secret:= "ABRACADABRA12345"
	key := []byte(secret)

	// encrypt value to base64
	cryptoText := Encrypt(key, originalText)

	// encrypt base64 crypto to original value
	text := Decrypt(key, cryptoText)
	if originalText != text {
		t.Error("It Doesn't work")
	}
}
