package tinyauth_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"fmt"
	"net/url"
	"strings"
	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
	"github.com/D10221/tinyauth/store"
)
func Test_Config(t *testing.T){
	// TODO:
	config:= config.NewConfig("1234567890ABCDEF")
	auth:= &tinyauth.TinyAuth{Config: config}
	if auth.Config.Secret == "" {
		t.Error("Config not Loaded")
	}

	if auth.Config.Secret != "1234567890ABCDEF" {
		t.Error("Config not Loaded")
	}

	auth = tinyauth.NewTinyAuth(config, nil)

	if auth.Config.Secret != "1234567890ABCDEF" {
		t.Error("Config not Loaded")
	}

	if e:= auth.Config.Validate() ; e!=nil {
		t.Errorf("Config not valid %s", e.Error())
	}

	if !auth.Config.Valid(){
		t.Error("Should be valid")
	}

	auth.Config.Secret = ""
	if auth.Config.Valid(){
		t.Error("Shouldnt be valid")
	}
}
func Test_RequireAuthentication_Encrypted(t *testing.T){

	config := &config.TinyAuthConfig{
		Secret: "ABRACADABRA12345",
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}

	tAuth:= tinyauth.NewTinyAuth(config, nil )

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

	tAuth.CredentialStore.Load( &store.Credential{"admin", pwd})

	found, e := tAuth.CredentialStore.FindByUserName("admin")
	if e!=nil && e!= store.NotFound || !found.Valid() {
		t.Error("Credential not Found")
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

	config := &config.TinyAuthConfig{
		Secret: "",
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}

	tauth:=  tinyauth.NewTinyAuth(config, nil )

	wrapped:= func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}

	handler:= tauth.RequireAuthentication(wrapped)

	request, _ := http.NewRequest("GET", "/", nil)

	response := httptest.NewRecorder()

	request.Header.Add(config.AuthorizationKey, tauth.Encoder.Encode("admin", "password"))

	tauth.CredentialStore.Load(&store.Credential{"admin", "password"})

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

func Test_GetFormCredentials(t *testing.T){

	//Configure Auth
	config := &config.TinyAuthConfig{
		Secret: "",
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}
	auth :=  tinyauth.NewTinyAuth(config, nil )
	auth.CredentialStore.Load(&store.Credential{"admin", "password"})

	// prep Request
	request, _ := http.NewRequest("GET", "/", nil)
	request.Form = url.Values{}
	request.Form.Add("username", "admin")
	request.Form.Add("password", "password")

	// set Writer
	response := httptest.NewRecorder()

	// Handler
	handler:= func(w http.ResponseWriter, r *http.Request) {
		cred, e := auth.GetFormCredentials(r)
		if e!=nil {
			http.Error(w, e.Error(), 500 )
			return
		}
		ok, e := auth.Authenticate(cred)
		if e!=nil{
			http.Error(w, e.Error(), 500 )
			return
		}
		if ok {
			fmt.Fprint(w, "ok")
			return
		}
		http.Error(w, "unauthorized", 401)
	}

	//Exec Handler
	handler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v:\n\tbody: %v", response.Code, response.Body)
	}

	//Test Result
	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		t.Error(err)
	}
	if string(body)!= "ok" {
		t.Error("bad body %s", string(body))
	}
}

func Test_GetFormCredentials_Encrypted(t *testing.T){

	//Configure Auth
	config := &config.TinyAuthConfig{
		Secret: "0123456789ABCDEF",
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}
	auth :=  tinyauth.NewTinyAuth(config, nil )
	password, e:= auth.Criptico.Encrypt("password")
	if e!=nil {
		t.Error(e)
		return
	}
	auth.CredentialStore.Load(&store.Credential{"admin", password})


	handler:= func(w http.ResponseWriter, r *http.Request) {
		cred, e := auth.GetFormCredentials(r)
		if e!=nil {
			http.Error(w, e.Error(), 500 )
			return
		}
		ok, e := auth.Authenticate(cred)
		if e!=nil {
			http.Error(w, e.Error(), 500 )
			return
		}
		if !ok {
			http.Error(w, "unauthorized", 401)
			return
		}

		fmt.Fprint(w, "ok")
		return
	}
	// ...
	{
		// Test 1 , good credentials
		request, _ := http.NewRequest("GET", "/", nil)
		request.Form = url.Values{}
		request.Form.Add("username", "admin")
		request.Form.Add("password", "password")
		writer := httptest.NewRecorder()
		handler(writer, request)
		testResponse:= ResponseTester(t)
		testResponse(writer, http.StatusOK, "ok")
	}
	{
		// Credentials Ok, PostForm, On Post
		// Test 1 , good credentials
		request, _ := http.NewRequest("POST", "/", nil)
		request.PostForm = url.Values{}
		request.PostForm.Add("username", "admin")
		request.PostForm.Add("password", "password")
		writer := httptest.NewRecorder()
		handler(writer, request)
		testResponse:= ResponseTester(t)
		testResponse(writer, http.StatusOK, "ok")
	}

	{
		// Credentials Ok, PostForm MISSING , On Post
		// Test 1 , good credentials
		request, _ := http.NewRequest("POST", "/", nil)
		request.Form = url.Values{}
		request.Form.Add("username", "admin")
		request.Form.Add("password", "password")
		writer := httptest.NewRecorder()
		handler(writer, request)
		testResponse:= ResponseTester(t)
		testResponse(writer, http.StatusOK, "ok")
	}

	{
		// Test2 Bad user
		request, _ := http.NewRequest("GET", "/", nil)
		request.Form = url.Values{}
		request.Form.Add("username", "username")
		request.Form.Add("password", "1234")
		writer := httptest.NewRecorder()
		handler(writer, request)
		expect := ResponseTester(t)
		expect(writer, http.StatusUnauthorized, "unauthorized")
	}

	{
		// Test3  bad password
		request, _ := http.NewRequest("GET", "/", nil)
		request.Form = url.Values{}
		request.Form.Add("username", "admin")
		request.Form.Add("password", "1234")
		writer := httptest.NewRecorder()
		handler(writer, request)
		expect := ResponseTester(t)
		expect(writer, http.StatusUnauthorized, "unauthorized")
	}
}

type ResponseTest  func(response *httptest.ResponseRecorder, expectedResponseCode int, expectedBody string );

func ResponseTester (t *testing.T) ResponseTest {

	return func(response *httptest.ResponseRecorder, expectedResponseCode int, expectedBody string ) {

		if response.Code != expectedResponseCode {
			t.Fatalf("Expected status code %v , received: %v :\n\t body: %v",expectedResponseCode, response.Code, response.Body)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil{
			t.Error(err)
		}
		content := string(body)
		if expectedBody != strings.TrimSpace(content){
			t.Errorf("Expected body %s \n received: %s", expectedBody ,content)
		}
	}
}

func Test_Auth_Encode(t *testing.T){
	auth:= tinyauth.NewTinyAuth(config.NewConfig("0123456789ABCDEF"), nil )
	key,value:= auth.Encode(&store.Credential{"admin", "password"})
	if key != auth.Config.AuthorizationKey {
		t.Error("Bad Key")
	}
	d,e:= auth.Encoder.Decode(value)
	if e!=nil {
		t.Error(e)
		return
	}
	if d != "admin:password" {
		t.Error("Bad Encoding")
	}
}


