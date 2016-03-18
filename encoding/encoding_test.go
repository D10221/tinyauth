package encoding

import (
	"testing"
	"strings"
)


func Test_Encoding(t *testing.T){
	encoded:= Encode("u", "p")

	scheme:= basicScheme+encoded

	decoded, err, ok:= Decode(scheme);
	if !ok {
		t.Errorf("Can't decode %s %s", encoded ,err.Error())
		return
	}
	if decoded != "u:p"{
		t.Errorf("Doesn't work, decoded: %s", decoded)
		return
	}
	parts := strings.Split(decoded, ":")
	if len(parts) < 2 {
		t.Error("Bad Decoding")
	}
	if parts[0] != "u" || parts[1] != "p"{
		t.Errorf("Bad decoding")
	}
}
