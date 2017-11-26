package isac

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	AccessToken       string `json:"AccessToken,omitempty"`
	AccessTokenSecret string `json:"AccessTokenSecret,omitempty"`
}

func NewConfig(configPath string) (config *Config, err error) {
	config = &Config{}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(b, config)
	return config, err
}
