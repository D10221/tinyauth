package tinyauth

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"fmt"
)

func Test_RequireAuthentication_Encrypted(t *testing.T){

	config := &TinyAuthConfig{
		Secret: "ABRACADABRA12345",
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}

	tAuth:= NewTinyAuth(config)

	wrapped:= func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}

	handler:= tAuth.RequireAuthentication(wrapped)

	request, _ := http.NewRequest("GET", "/", nil)

	response := httptest.NewRecorder()

	// encoded,  not encrypted password over the wire : https required
	request.Header.Add("Authorization", tAuth.Encoder.Encode("admin", "password"))

	pwd, e:= tAuth.Criptico.Encrypt("password")

	if e!=nil {
		t.Error(e)
	}

	tAuth.CredentialStore.Load( Credential{"admin", pwd})

	found, err := tAuth.CredentialStore.FindUser("admin")
	if err!= nil || !found.Valid() {
		t.Error("user not Found")
	}

	handler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v:\n\tbody: %v", response.Code, response.Body)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		t.Error(err)
	}
	if string(body)!= "ok" {
		t.Error("bad body %s", string(body))
	}
}

func Test_RequireAuthentication(t *testing.T){

	config := &TinyAuthConfig{
		Secret: "",
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}

	tauth:=  NewTinyAuth(config)

	wrapped:= func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}

	handler:= tauth.RequireAuthentication(wrapped)

	request, _ := http.NewRequest("GET", "/", nil)

	response := httptest.NewRecorder()

	request.Header.Add(config.AuthorizationKey, tauth.Encoder.Encode("admin", "password"))

	tauth.CredentialStore.Load(Credential{"admin", "password"})

	handler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v:\n\tbody: %v", response.Code, response.Body)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		t.Error(err)
	}
	if string(body)!= "ok" {
		t.Error("bad body %s", string(body))
	}
}


