package credentials

import "testing"

func Test_Credentials(t *testing.T){

	credentials:= New("bob", "P@55w0rd!")

	encoded:= credentials.Encode()

	decodedCredentials, err:= Decode(encoded)

	if err!= nil  {
		t.Error(err)
	}

	if decodedCredentials.Username != "bob" || decodedCredentials.Password!= "P@55w0rd!"{
		t.Error("It Doesn't Work")
	}

	if !credentials.Authenticate() {
		t.Error("anauthorized")
	}
}
