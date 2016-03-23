package encoder

import (
	"encoding/base64"
	"strings"
)

type Encoder interface {
	Encode(left, right  string) string
	Decode(auth string) (decoded string, err error)
}

type DefaultEncoder struct {
	BasicScheme string
}

func (enc *DefaultEncoder) Encode(left, right string) string {
	return enc.BasicScheme + base64.StdEncoding.EncodeToString([]byte(left + ":" + right))
}

func (enc *DefaultEncoder) Decode(auth string) (decoded string, err error) {
	s := enc.skipScheme(auth)
	b, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return "", e
	}
	return string(b), nil
}


func (encoder *DefaultEncoder) skipScheme(auth string) (string) {
	if strings.HasPrefix(auth, encoder.BasicScheme) {
		return auth[len(encoder.BasicScheme):]
	}
	return auth
}


