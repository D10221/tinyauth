// tinyauth configuration
package config

import (
	"io/ioutil"
	"encoding/json"
	"errors"
)

type Config struct {
	Secret string
	AuthorizationKey string
}

var Current = &Config{}

// Load config From Json file (full path)
func (config *Config) LoadConfig(path string) error {

	bytes, err :=  ioutil.ReadFile(path)
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
func (config *Config)Validate() error {
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