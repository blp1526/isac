package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	AccessToken       string `json:"AccessToken,omitempty"`
	AccessTokenSecret string `json:"AccessTokenSecret,omitempty"`
	Zone              string `json:"Zone,omitempty"`
}

func New(configPath string) (config *Config, err error) {
	config = &Config{}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(b, config)
	return config, err
}

func (config *Config) CreateFile(dir string) (err error) {
	_, err = os.Stat(dir)
	if err == nil {
		return fmt.Errorf("Already you have %s", dir)
	}

	j, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dir, j, 0600)
	if err != nil {
		return err
	}

	return nil
}
