package isac

import (
	"github.com/sirupsen/logrus"
)

type Isac struct {
	configPath string
	logger     *logrus.Logger
	zones      string
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
		zones:      zones,
	}
	return i
}

func (i *Isac) Run() (err error) {
	i.logger.Debugf("configPath: %s", i.configPath)
	i.logger.Debugf("zones: %s", i.zones)

	config, err := loadConfig(i.configPath)
	if err != nil {
		return err
	}

	i.logger.Debugf("AccessToken: %s", config.AccessToken)
	i.logger.Debugf("AccessTokenSecret: %s", config.AccessTokenSecret)
	return nil
}
