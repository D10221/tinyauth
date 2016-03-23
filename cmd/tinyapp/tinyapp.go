package tinyapp

import (
	"github.com/D10221/tinyauth"
	"net/http"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"runtime"
	"github.com/D10221/tinyauth/store"
)

type TinyApp struct {
	Auth       *tinyauth.TinyAuth
	currentDir string
	Templates  string
}

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

func trace() {
	pc := make([]uintptr, 1)  // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
}

func (app *TinyApp) Authenticate(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	if e:= r.ParseForm(); e!=nil {
		log.Printf("Error: %s", e.Error())
		http.Error(w, e.Error() , 500)
	}

	credential, err := app.Auth.GetFormCredentials(r)

	if !credential.Valid(){
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		log.Printf("No Credentials")
		return
	}
	log.Print(credential)

	found := app.Auth.CredentialStore.FindUser(credential.Username)
	log.Printf("Found: %v", found)

	ok, err:= app.Auth.Authenticate(credential)
	if err!=nil {
		log.Printf("Error: %s", err.Error())
		http.Error(w, err.Error(), 500 )
		return
	}
	if !ok {
		log.Printf("Bad Credentials")
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "%s: Ok", credential.Username)
	log.Printf("%s: Ok", credential.Username)
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

