package tinyauth

import "testing"

func Test_Config_Valid(t *testing.T){
	config:= TinyAuthConfig{}
	if e:= config.Validate() ; e==nil {
		t.Error("It Shoudln't")
	}
	config = TinyAuthConfig{BasicScheme:"", AuthorizationKey:"", Secret:""}
	if e:= config.Validate() ; e==nil {
		t.Error("It Shoudln't")
	}
}