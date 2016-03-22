package tinyauth

import (
	"encoding/base64"
	"strings"
)

type Encoder interface {
	Encode(left, right  string) string
	Decode(auth string) (decoded string, err error)
	EncodeWithSchema(left, right string) string
}


type DefaultEncoder struct {
	config *TinyAuthConfig
}

func (enc *DefaultEncoder) Encode(left, right  string) string {
	return base64.StdEncoding.EncodeToString([]byte(left + ":" + right))
}
func (enc *DefaultEncoder) EncodeWithSchema(left, right string) string {
	return enc.config.BasicScheme + enc.Encode(left, right)
}

func (enc *DefaultEncoder) Decode(auth string) (decoded string, err error) {
	s, e := skipScheme(enc.config, auth)
	if e != nil {
		return "", e
	}
	b, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return "", e
	}
	return string(b), nil
}

func (enc *DefaultEncoder) ShouldDecode(auth string) string {
	r, err := enc.Decode(auth)
	if err != nil {
		return ""
	}
	return r
}

type DecodeError struct {
	message string
}

func (e *DecodeError) Error() string {
	return e.message
}


func skipScheme(config *TinyAuthConfig, auth string) (string, error) {
	if strings.HasPrefix(auth, config.BasicScheme) {
		return auth[len(config.BasicScheme):], nil
	}
	return "", &DecodeError{"No Scheme"}
}
