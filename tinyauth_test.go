package tinyauth

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"fmt"
	"github.com/D10221/tinyauth/credentials"
	"github.com/D10221/tinyauth/config"
)

func Test_RequireAuthentication(t *testing.T){

	wrapped:= func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}

	handler:= RequireAuthentication(wrapped)

	request, _ := http.NewRequest("GET", "/", nil)

	response := httptest.NewRecorder()

	request.Header.Add("Authorization", "Basic YWRtaW46UEA1NXcwcmQh")

	credentials.Credentials.Load(credentials.Credential{"admin", "P@55w0rd!"})

	if(credentials.Credentials.All()[0].Username != "admin"){
		t.Error("no credentials")
	}

	config.Current.AuthorizationKey = "Authorization"

	handler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %v", "200", response.Code)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		t.Error(err)
	}
	if string(body)!= "ok" {
		t.Error("bad body %s", string(body))
	}
}

