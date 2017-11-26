package isac

import (
	"strings"

	"github.com/blp1526/isac/lib/server"
	"github.com/sirupsen/logrus"
)

type Isac struct {
	configPath string
	logger     *logrus.Logger
	zones      []string
}

func New(configPath string, verbose bool, zones string) *Isac {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}

	logger := logrus.New()
	logger.Formatter = formatter
	if verbose {
		logger.Level = logrus.DebugLevel
	}

	i := &Isac{
		configPath: configPath,
		logger:     logger,
		zones:      strings.Split(zones, ","),
	}
	return i
}

func (i *Isac) Run() (err error) {
	i.logger.Debugf("configPath: %s", i.configPath)
	i.logger.Debugf("zones: %v", i.zones)

	config, err := NewConfig(i.configPath)
	if err != nil {
		return err
	}

	i.logger.Debugf("AccessToken: %s", config.AccessToken)
	i.logger.Debugf("AccessTokenSecret: %s", config.AccessTokenSecret)

	servers, err := server.All(i.zones)
	if err != nil {
		return err
	}

	i.logger.Debugf("servers: %v", servers)
	return nil
}
