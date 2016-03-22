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

	request.Header.Add("Authorization", "Basic YWRtaW46UEA1NXcwcmQh")

	tAuth.CredentialStore.Load( Credential{"admin", "P@55w0rd!"})

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
		Secret: "ABRACADABRA12345",
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

	request.Header.Add("Authorization", "Basic YWRtaW46Y0p0UFFLYzJ2N0xNYkFNT1RWQmhvMHU2bXVzPQ==")

	tauth.CredentialStore.Load(Credential{"admin", "password"})

	if tauth.CredentialStore.All()[0].Username != "admin" {
		t.Error("There's something wrong with Credential Store implementation")
	}

	found, err := tauth.CredentialStore.FindUser("admin")
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


