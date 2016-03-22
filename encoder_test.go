package tinyauth

import (
	"testing"
	"strings"
)

func Test_Encoder(t *testing.T){

	var encoder Encoder  = &DefaultEncoder{BasicScheme: "Basic "}

	encoded:= encoder.Encode("left", "right")

	if !strings.HasPrefix(encoded, "Basic "){
		t.Error("Missing Scheme")
	}

	out,e := encoder.Decode(encoded)

	if e!=nil {
		t.Error(e)
	}

	if out!= "left:right" {
		t.Errorf("Fail: %s", out )
	}
}
