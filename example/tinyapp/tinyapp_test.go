package tinyapp

import (
	"testing"
	"log"
	"net/url"
	"net/http/httptest"

	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
	"net/http"
	"github.com/D10221/tinyauth/store"
)

func Test_TinyApp(t *testing.T){

	app:=  &TinyApp{Auth: tinyauth.NewTinyAuth( config.NewConfig("0123456789ABCDEF"))}

	if app.Auth.Config.Secret == "" {
		t.Error("No Secret")
	}

	credential, e := app.Auth.NewCredential("user", "password")
	if e!=nil {
		t.Error(e)
		return;
	}
	if !credential.Valid() {
		t.Error("not valid credential")
	}
	if credential.Username != "user" {
		t.Error("Bad credential")
	}

	if credential.Password == "password" {
		t.Error("credential is not encrypted")
	}
}

func Test_TinyApp_Config(t *testing.T){

	app:=  &TinyApp{Auth: tinyauth.NewTinyAuth( &config.TinyAuthConfig{})}

	log.Printf("App Dir: %s", app.CurrentDir())

	path:= app.MakePath("config.json")

	log.Printf("App config: %s", path)

	e:= app.Auth.LoadConfig(path)

	if e!=nil {
		t.Error(e)
		return
	}

	if app.Auth.Config.Secret == "" {
		t.Error("No Secret")
		return
	}

	credential, e := app.Auth.NewCredential("user", "password")

	if e!=nil {
		t.Error(e)
		return
	}
	if !credential.Valid() {
		t.Error("not valid credential")
	}

	if credential.Username != "user" {
		t.Error("Bad credential")
	}

	if credential.Password == "password" {
		t.Error("credential is not encrypted")
	}


}

func Test_TinyApp_NoEncryption(t *testing.T){

	app:=  &TinyApp{Auth: tinyauth.NewTinyAuth( &config.TinyAuthConfig{})}

	if app.Auth.Config.Secret != "" {
		t.Error("Should not have Secret")
	}

	credential, e := app.Auth.NewCredential("user", "password")
	if e!= nil {
		t.Error(e)
		t.Fail()
		return
	}

	if !credential.Valid() {
		t.Error("not valid credential")
	}
	if credential.Username != "user" {
		t.Error("Bad credential user")
	}

	if credential.Password != "password" {
		t.Error("Bad credential password")
	}
	log.Println(credential)
}


func makeRequest(method string, username string,password string) *http.Request {
	request, _ := http.NewRequest(method, "/", nil)
	if method == http.MethodPost {
		request.Form = url.Values{}
		request.Form.Add("username", username)
		request.Form.Add("password", password)
		request.PostForm = url.Values{}
		copyValues(request.PostForm, request.Form)
	} else if method == http.MethodGet {
		request.Header.Add("username", username)
		request.Header.Add("password", username)
	}

	return request
}
func copyValues(dst, src url.Values) {
	for k, vs := range src {
		for _, value := range vs {
			dst.Add(k, value)
		}
	}
}

func Test_Authenticate(t *testing.T){

	app:= &TinyApp{Auth: tinyauth.NewTinyAuth(config.NewConfig("0123456789ABCDEF"))}
	app.Auth.CredentialStore.Add(&store.Credential{"admin", "password"})
	app.Auth.CredentialStore.UpdateWhere(store.ByName("admin"), app.Auth.EncryptPassword )
	ok, err:= app.Auth.Authenticate(&store.Credential{"admin", "password"})

	if err!=nil { t.Error(err)  ; return }

	if !ok { t.Error("Authentication system failure") ; return }

	// Get
	{
		// Setup ...
		writer := httptest.NewRecorder()
		r, _:= http.NewRequest(http.MethodGet, "/", nil )
		key, value:= app.Auth.Encode(&store.Credential{"admin", "password"})
		r.Header.Add(key, value)
		// Test subject ...
		app.Authenticate(writer, r)
		// Test ...
		if writer.Code  !=  http.StatusOK {
			t.Log(string(writer.Body.Bytes()))
			t.Errorf("method: %s Expected %d , got: %d", r.Method , http.StatusOK, writer.Code)
		}
	}
	// Post
	{
		// Setup ...
		writer := httptest.NewRecorder()
		r:= makeRequest(http.MethodPost,"admin", "password")
		// Test subject ...
		app.Authenticate(writer, r)
		// Test ...
		if writer.Code  != http.StatusOK {
			t.Log(string(writer.Body.Bytes()))
			t.Errorf("method: %s Expected %d , got: %d", r.Method , http.StatusOK, writer.Code)
		}
	}
	// Not Get Nor Post
	{
		//Setup ...
		writer := httptest.NewRecorder()
		r:= makeRequest(http.MethodPut,"admin", "password")
		// Test subject ...
		app.Authenticate(writer, r)
		// Test ...
		if writer.Code  != http.StatusMethodNotAllowed {
			t.Errorf("method: %s Expected %d , got: %d", r.Method , http.StatusMethodNotAllowed, writer.Code)
		}
	}
}





