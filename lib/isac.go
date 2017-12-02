package isac

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blp1526/isac/lib/api"
	"github.com/blp1526/isac/lib/config"
	"github.com/blp1526/isac/lib/resource/server"
	"github.com/blp1526/isac/lib/row"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
)

type Isac struct {
	client       *api.Client
	config       *config.Config
	showServerID bool
	logger       *logrus.Logger
	row          *row.Row
	servers      []server.Server
	zones        []string
}

func New(configPath string, showServerID bool, verbose bool, zones string) (i *Isac, err error) {
	config, err := config.New(configPath)
	if err != nil {
		return i, err
	}

	client := api.NewClient(config.AccessToken, config.AccessTokenSecret)
	row := row.New()

	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}

	logger := logrus.New()
	logger.Formatter = formatter
	if verbose {
		logger.Level = logrus.DebugLevel
	}

	i = &Isac{
		client:       client,
		config:       config,
		showServerID: showServerID,
		logger:       logger,
		row:          row,
		zones:        strings.Split(zones, ","),
	}
	return i, nil
}

func (i *Isac) Run() (err error) {
	err = i.reloadServers()
	if err != nil {
		return err
	}

	err = termbox.Init()
	if err != nil {
		return err
	}

	defer termbox.Close()

	i.draw()

MAINLOOP:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break MAINLOOP
			}
		default:
			i.draw()
		}
	}
	return nil
}

func (i *Isac) setLine(y int, line string) {
	runes := []rune(line)
	x := 0
	for _, r := range runes {
		bgColor := termbox.ColorDefault
		termbox.SetCell(x, y, r, termbox.ColorWhite, bgColor)
		x += runewidth.RuneWidth(r)
	}
}

func (i *Isac) draw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	for index, header := range i.row.Headers() {
		i.setLine(index, header)
	}

	offsetSize := len(i.row.Headers())
	for index, server := range i.servers {
		i.setLine(index+offsetSize, server.String(i.showServerID))
	}
	termbox.Flush()
}

func (i *Isac) reloadServers() (err error) {
	for _, zone := range i.zones {
		statusCode, respBody, err := i.client.Request("GET", zone, "server", "", nil)
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

		i.logger.Debugf("zone: %v, Servers.Count: %v", zone, sc.Count)

		for _, s := range sc.Servers {
			i.servers = append(i.servers, s)
		}
	}

	return nil
}
