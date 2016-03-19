package config

import (
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Secret string
	AuthorizationKey string
}

var Current = &Config{}

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