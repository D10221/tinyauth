package config

import (
	"testing"
	"os"
	"path/filepath"
)

func Test_Config_Valid(t *testing.T){

	config:= TinyAuthConfig{}

	if e:= config.Validate() ; e==nil {
		t.Error("It Shoudln't")
	}
	config = TinyAuthConfig{BasicScheme:"", AuthorizationKey:"", Secret:""}
	if e:= config.Validate() ; e==nil {
		t.Error("It Shoudln't")
	}

	dir , e := os.Getwd()
	if e!= nil {
		t.Error(e)
	}

	e= config.LoadConfig(filepath.Join(dir, "testdata/config.json"))

	if e!=nil {
		t.Error(e)
	}

	if(config.Secret == "") {
		t.Error("config not loaded")
	}

	if(config.Secret != "ABRACADABRA12345") {
		t.Error("config not loaded")
	}

}