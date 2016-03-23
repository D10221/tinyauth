package config

import (
	"io/ioutil"
	"encoding/json"
	"errors"
)

type TinyAuthConfig struct {
	Secret           string
	AuthorizationKey string
	BasicScheme      string // "Basic "
}

// Load config From Json file (full path)
func (config *TinyAuthConfig) LoadConfig(path string) error {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}

// Is configuration valid ?
// maybe
func (config *TinyAuthConfig)Validate() error {
	if config.AuthorizationKey == "" {
		return errors.New("No Authorization Key")
	}
	if config.Secret == "" {
		return errors.New("No secret")
	}

	if len(config.Secret) != 16 {
		return errors.New("Bad secret length")
	}

	return nil
}

func (config *TinyAuthConfig)Valid() bool {
	e:= config.Validate()
	return e == nil
}

func NewConfig(secret string) *TinyAuthConfig {
	return &TinyAuthConfig{
		Secret: secret,
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	}
}

