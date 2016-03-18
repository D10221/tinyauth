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

func Decode(auth string) (decoded string, err error, ok bool) {
	if s, e, ok := skipScheme(auth); !ok{
		return "", e, false
	} else {
		b, e := base64.StdEncoding.DecodeString(s)
		if e != nil {
			return "", e, false
		}
		return string(b), nil , true
	}
}
func skipScheme(auth string)( string , error , bool ){
	if strings.HasPrefix(auth, basicScheme)	{
		return auth[len(basicScheme):] , nil , true
	}
	return "", &DecodeError{"No Scheme"}, false
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