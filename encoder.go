package tinyauth

import (
	"encoding/base64"
	"strings"
	"errors"
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
	s, e := enc.skipScheme(auth)
	if e != nil {
		return "", e
	}
	b, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return "", e
	}
	return string(b), nil
}


func (encoder *DefaultEncoder) skipScheme(auth string) (string, error) {
	if strings.HasPrefix(auth, encoder.BasicScheme) {
		return auth[len(encoder.BasicScheme):], nil
	}
	return "", errors.New("No Scheme")
}
