package isac

import (
	"encoding/json"
	"errors"
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
	i.draw("OK")

MAINLOOP:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break MAINLOOP
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				i.currentRowUp()
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				i.currentRowDown()
			}
		default:
			i.draw("OK")
		}
	}
	return nil
}

func (i *Isac) setLine(y int, line string) {
	runes := []rune(line)
	x := 0
	for _, r := range runes {
		fgColor := termbox.ColorDefault
		bgColor := termbox.ColorDefault

		if i.row.Current == y {
			fgColor = termbox.ColorBlack
			bgColor = termbox.ColorYellow
		}

		termbox.SetCell(x, y, r, fgColor, bgColor)
		x += runewidth.RuneWidth(r)
	}
}

func (i *Isac) draw(status string) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	headers := i.row.Headers(status)

	if i.row.Current == 0 {
		i.row.Current = len(headers)
	}

	for index, header := range headers {
		i.setLine(index, header)
	}

	offsetSize := len(headers)
	i.row.MovableTop = len(headers)

	for index, server := range i.servers {
		no := index + 1
		server.No = no
		i.setLine(index+offsetSize, server.String(i.showServerID))
	}
	i.row.MovableBottom = len(i.servers) + len(headers) - 1

	termbox.Flush()
}

func (i *Isac) currentRowUp() {
	if i.row.Current > i.row.MovableTop {
		i.row.Current -= 1
	}

	i.draw("OK")
}

func (i *Isac) currentRowDown() {
	if i.row.Current < i.row.MovableBottom {
		i.row.Current += 1
	}

	i.draw("OK")
}

func (i *Isac) reloadServers() (err error) {
	for _, zone := range i.zones {
		statusCode, respBody, err := i.client.Request("GET", zone, []string{"server"}, nil)
		if err != nil {
			return err
		}

		if statusCode != 200 {
			return errors.New(fmt.Sprintf("statusCode: %v", statusCode))
		}

		sc := server.NewCollection(zone)
		err = json.Unmarshal(respBody, sc)
		if err != nil {
			return err
		}

		for _, s := range sc.Servers {
			i.servers = append(i.servers, s)
		}
	}

	return nil
}
