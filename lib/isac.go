package isac

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blp1526/isac/lib/api"
	"github.com/blp1526/isac/lib/config"
	"github.com/blp1526/isac/lib/resource/server"
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

	c, err := config.New(i.configPath)
	if err != nil {
		return err
	}

	i.logger.Debugf("AccessToken: %s", c.AccessToken)
	i.logger.Debugf("AccessTokenSecret: %s", c.AccessTokenSecret)

	ac := api.NewClient(c.AccessToken, c.AccessTokenSecret)

	servers := []server.Server{}
	for _, zone := range i.zones {
		statusCode, respBody, err := ac.Request("GET", zone, "server", "", nil)
		if err != nil {
			return err
		}

		if statusCode != 200 {
			return fmt.Errorf("statusCode: %v", statusCode)
		}

		sc := server.NewCollection()
		err = json.Unmarshal(respBody, sc)
		if err != nil {
			return err
		}

		i.logger.Debugf("sc.Count: %v", sc.Count)

		for _, s := range sc.Servers {
			servers = append(servers, s)
		}

	}

	for _, server := range servers {
		i.logger.Debugf("server: %+v", server)
	}

	return nil
}
