package isac

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	AccessToken       string `json:"AccessToken,omitempty"`
	AccessTokenSecret string `json:"AccessTokenSecret,omitempty"`
}

func loadConfig(path string) (config *Config, err error) {
	config = &Config{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(b, config)
	return config, err
}
