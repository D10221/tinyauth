package encoding

import (
	"encoding/base64"
	"strings"
)


func Encode(left, right  string) string {
	return base64.StdEncoding.EncodeToString([]byte(left + ":" + right))
}
func EncodeWithSchema(left, right string) string {
	return basicScheme+ Encode(left, right)
}

var basicScheme string

func Decode(auth string) (decoded string, err error) {
	s, e := skipScheme(auth)
	if e!=nil {
		return "", e
	}
	b, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return "", e
	}
	return string(b), nil
}

func ShouldDecode(auth string) string {
	r,err := Decode(auth)
	if err!=nil {
		return ""
	}
	return r
}

func skipScheme(auth string)( string , error  ){
	if strings.HasPrefix(auth, basicScheme)	{
		return auth[len(basicScheme):] , nil
	}
	return "", &DecodeError{"No Scheme"}
}

type DecodeError struct {
	message string
}

func(e *DecodeError) Error() string{
	return e.message
}

func init(){
	basicScheme = "Basic "
}