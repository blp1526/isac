package isac

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/blp1526/isac/lib/api"
	"github.com/blp1526/isac/lib/config"
	"github.com/blp1526/isac/lib/resource/server"
	"github.com/blp1526/isac/lib/row"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
)

type Isac struct {
	filter       string
	client       *api.Client
	config       *config.Config
	showServerID bool
	logger       *logrus.Logger
	row          *row.Row
	// FIXME: wasteful
	serverByCurrentRow map[int]server.Server
	servers            []server.Server
	reverseSort        bool
	zones              []string
	message            string
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

	zs := []string{config.Zone}
	if zones != "" {
		zs = strings.Split(zones, ",")
	}

	i = &Isac{
		client:       client,
		config:       config,
		showServerID: showServerID,
		logger:       logger,
		row:          row,
		zones:        zs,
		reverseSort:  false,
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
			case termbox.KeyCtrlU:
				message := i.currentServerUp()
				i.draw(message)
			case termbox.KeyCtrlR:
				i.refresh()
			case termbox.KeyBackspace2, termbox.KeyCtrlH:
				i.removeRuneFromFilter()
			case termbox.KeyCtrlS:
				i.reverseSort = !i.reverseSort
				i.draw("")
			default:
				if ev.Ch != 0 {
					i.addRuneToFilter(ev.Ch)
				}
			}
		default:
			i.draw("")
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

func (i *Isac) draw(message string) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	if message != "" {
		i.message = message
	}

	var servers []server.Server
	for _, s := range i.servers {
		if strings.Contains(s.Name, i.filter) {
			servers = append(servers, s)
		}
	}

	sort.Slice(servers, func(x, y int) bool {
		if i.reverseSort {
			return servers[x].Instance.Status > servers[y].Instance.Status
		} else {
			return servers[x].Instance.Status < servers[y].Instance.Status
		}
	})

	if i.row.Current == 0 {
		i.row.Current = i.row.HeadersSize()
	}

	i.row.MovableBottom = len(servers) + i.row.HeadersSize() - 1
	if i.row.Current > i.row.MovableBottom {
		i.row.Current = i.row.HeadersSize()
	}

	headers := i.row.Headers(i.message, strings.Join(i.zones, ", "), len(servers), i.currentNo(), i.filter)

	for index, header := range headers {
		i.setLine(index, header)
	}

	i.serverByCurrentRow = map[int]server.Server{}

	for index, server := range servers {
		currentRow := index + i.row.HeadersSize()
		i.setLine(currentRow, server.String(i.showServerID))
		i.serverByCurrentRow[currentRow] = server
	}

	termbox.Flush()
}

func (i *Isac) currentRowUp() {
	if i.row.Current > i.row.HeadersSize() {
		i.row.Current -= 1
	}

	i.draw("")
}

func (i *Isac) currentRowDown() {
	if i.row.Current < i.row.MovableBottom {
		i.row.Current += 1
	}

	i.draw("")
}

func (i *Isac) reloadServers() (err error) {
	i.servers = []server.Server{}

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

func (i *Isac) currentNo() int {
	return i.row.Current + 1 - i.row.HeadersSize()
}

func (i *Isac) currentServerUp() (message string) {
	s := i.serverByCurrentRow[i.row.Current]

	if s.ID == "" {
		return "[ERROR] Current row has no Server"
	}

	if s.Instance.Status == "up" {
		return fmt.Sprintf("Server.Name %v is already up", s.Name)

	}

	paths := []string{"server", s.ID, "power"}
	statusCode, _, err := i.client.Request("PUT", s.Zone.Name, paths, nil)

	if err != nil {
		return fmt.Sprintf("[ERROR] %v", err)
	}

	if statusCode != 200 {
		return fmt.Sprintf("[ERROR] Server.Name: %v, PUT /server/:id/power failed, statusCode: %v", s.Name, statusCode)
	}

	return fmt.Sprintf("Server.Name %v is booting, wait few seconds, and refresh", s.Name)
}

func (i *Isac) refresh() {
	var message string

	err := i.reloadServers()
	if err != nil {
		message = fmt.Sprintf("[ERROR] %v", err)
	}

	if message == "" {
		message = "Servers have been refreshed"
	}

	i.draw(message)
}

func (i *Isac) addRuneToFilter(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	i.filter = i.filter + string(buf[:n])
	i.draw("")
}

func (i *Isac) removeRuneFromFilter() {
	r := []rune(i.filter)
	if len(r) > 0 {
		i.filter = string(r[:(len(r) - 1)])
	}

	i.draw("")
}
