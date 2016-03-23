package main

import (
	"testing"
	"os"
	"strings"
	"bytes"
	"fmt"
	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
)

func Test_TinyAUthTool_NoTask(t *testing.T){
	args:= make([]string, 0)
	// append(args, "--config:=")
	// args:= [...]string { "-task=encode", "-left=me" , "-right=me" , "-config=.noconfig.json"}
	var app = &TinyApp{Auth : tinyauth.NewTinyAuth(&config.TinyAuthConfig{
		Secret: "", // ABRACADABRA12345
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	})}
	e:= app.ParseCommandLine(args)
	if e != nil {
		t.Errorf("it shouldn't fail %s", e.Error())
	}
	e = app.SwitchTask(os.Stdout)
	if message := func() string{ if e!= nil { return e.Error() } else { return ""  }}(); message != "task required" {
		t.Errorf("Error: %s \n", message)
	}
}

func Test_TinyAUthTool(t *testing.T){

	var app = &TinyApp{Auth : tinyauth.NewTinyAuth(&config.TinyAuthConfig{
		Secret: "", // ABRACADABRA12345
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	})}
	args:= [...]string { "-task=encode", "-left=me" , "-right=me"}
	e:= app.ParseCommandLine(args[:])
	if e != nil {
		t.Error("it shouldn't fail")
	}
	writer:= &bytes.Buffer{}
	e = app.SwitchTask(writer)
	if e != nil {
		t.Error("it shouldn't fail")
	}
	result:= string(writer.Bytes())
	if !strings.Contains(result, "Basic "){
		t.Errorf("Nope: %v", result)
	}
}

func Test_TinyAUthLoadConfig(t *testing.T){

	var app = &TinyApp{Auth : tinyauth.NewTinyAuth(&config.TinyAuthConfig{
		Secret: "", // ABRACADABRA12345
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	})}

	args:= [...]string { "-task=encode", "-left=me" , "-right=me" , "-config=config.json"}
	e:= app.ParseCommandLine(args[:])
	if e != nil {
		t.Error("it shouldn't fail")
	}
	writer:= &bytes.Buffer{}
	e = app.SwitchTask(writer)
	if e != nil {
		t.Error("it shouldn't fail")
	}
	result:= string(writer.Bytes())
	if !strings.Contains(result, "Basic "){
		t.Errorf("Nope: %v", result)
	}
}
func Test_TinyauthCantFindConfig(t *testing.T){

	args:= [...]string { "-task=encode", "-left=me" , "-right=me" , "-config=.noconfig.json"}
	var app = &TinyApp{Auth : tinyauth.NewTinyAuth(&config.TinyAuthConfig{
		Secret: "", // ABRACADABRA12345
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	})}

	e:= app.ParseCommandLine(args[:])
	if e != nil {
		t.Error(e)
	}

	e = app.LoadConfig()
	if e!= nil {
		return
	}
	code:= BonkCode(e)
	if code!= 404 {
		t.Errorf("Bad Error Code: %s", code)
	}
	t.Logf("Success %s", e.Error())
}

func Test_BonkCode(t *testing.T){
	bonk:= &Bonk{"x", 500}
	code := BonkCode(bonk)
	if code != 500 {
		t.Error("Failed")
	}
}

func Test_TinyAUthLoadBadConfig(t *testing.T){

	var app = &TinyApp{Auth : tinyauth.NewTinyAuth(&config.TinyAuthConfig{
		Secret: "", // ABRACADABRA12345
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	})}

	args:= [...]string { "-task=encode", "-left=me" , "-right=me" , "-config=../../testdata/badconfig.json"}
	e:= app.ParseCommandLine(args[:])
	if e != nil {
		t.Errorf("it shouldn't fail:... %s", e.Error())
	}
	writer:= &bytes.Buffer{}
	e = app.SwitchTask(writer)
	if e != nil {
		t.Error("it shouldn't fail")
	}
	result:= string(writer.Bytes())
	if !strings.Contains(result, "Basic "){
		t.Errorf("Nope: %v", result)
	}
}

func Test_TestWriter(t *testing.T){
	writer:= &bytes.Buffer{}
	fmt.Fprintf(writer, "%v", 1)
	result := string(writer.Bytes())
	if result != "1" {
		t.Errorf("It Doesn't work: %v", result)
	}
}

