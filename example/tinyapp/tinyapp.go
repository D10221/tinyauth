package tinyapp

import (
	"github.com/D10221/tinyauth"
	"net/http"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"github.com/D10221/tinyauth/store"
)

type TinyApp struct {
	Auth       *tinyauth.TinyAuth
	currentDir string
	Templates  string
}

// A view:
func (app *TinyApp) Login(w http.ResponseWriter, r *http.Request) {

	t, e := template.ParseFiles(app.TemplatePath("login.html"))

	if e != nil {
		http.Error(w, e.Error(), 404)
		return
	}

	credential, _ := app.Auth.GetRequestCredentials(r)

	data := &struct {
		Title string
		FormAction string
		Credential *store.Credential
	}{
		Title: "Login Form",
		FormAction: "/authenticate",
		Credential: credential,
	}

	log.Printf("Header credentials %v ", data.Credential)
	log.Printf("Action: %v ", data.FormAction)

	e = t.Execute(w, data)

	if e != nil {
		http.Error(w, e.Error(), 500)
	}
}

func (app *TinyApp) CurrentDir() string {
	if app.currentDir != "" {
		return app.currentDir
	}
	dir, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	app.currentDir = dir
	return dir
}

// Nice: thanks to : http://stackoverflow.com/a/25927915/1901532
// needs import "runtime"
/*func trace() {
	pc := make([]uintptr, 1)  // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
}*/

func (app *TinyApp) Authenticate(w http.ResponseWriter, r *http.Request) {

	credential:= &store.Credential{}

	if r.Method == "POST" {
		c,e := app.Auth.GetFormCredentials(r)
		if  e != nil && e!= tinyauth.FormHasNoCredentials{
			http.Error(w, e.Error() , 500)
			return
		}
		credential = c

	} else if r.Method == "GET" {
		c, e := app.Auth.GetRequestCredentials(r)
		if  e != nil {
			http.Error(w, e.Error() , 500)
			return
		}
		credential = c
	} else {
		http.Error(w, "MethodNotAllowed" ,http.StatusMethodNotAllowed)
	}

	ok, err:= app.Auth.Authenticate(credential)

	if err!=nil {
		http.Error(w, err.Error(), 401 )
		return
	}

	if !ok {
		http.Error(w, tinyauth.UnAuthorized.Error(), 401 )
		return
	}

	fmt.Fprintf(w, "%s: Ok", credential.Username)
}

func( app *TinyApp) Secret(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Secret: %s ", "hello")
}

func (app *TinyApp) MakePath(path string) string {
	return filepath.Join(app.CurrentDir(), path)
}

func (app *TinyApp) TemplatePath(path string) string {
	return filepath.Join(app.CurrentDir(), app.Templates, path)
}

func (app *TinyApp) EncryptPassword (in *store.Credential) (*store.Credential, error ){
	return  app.Auth.EncryptPassword(in)
}

// ...

type Token struct {

}