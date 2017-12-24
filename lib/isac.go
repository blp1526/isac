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
)

type Isac struct {
	filter      string
	client      *api.Client
	config      *config.Config
	detail      bool
	unanonymize bool
	row         *row.Row
	// FIXME: wasteful
	serverByCurrentRow map[int]server.Server
	servers            []server.Server
	reverseSort        bool
	zones              []string
	message            string
}

func New(configPath string, unanonymize bool, zones string) (i *Isac, err error) {
	config, err := config.New(configPath)
	if err != nil {
		return i, err
	}

	client := api.NewClient(config.AccessToken, config.AccessTokenSecret)
	row := row.New()

	zs := []string{config.Zone}
	if zones != "" {
		zs = strings.Split(zones, ",")
	}

	i = &Isac{
		client:      client,
		config:      config,
		detail:      false,
		unanonymize: unanonymize,
		row:         row,
		zones:       zs,
		reverseSort: false,
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
			case termbox.KeyEnter:
				i.detail = !i.detail
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

		if !i.detail && i.row.Current == y {
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

	if i.detail {
		server := i.serverByCurrentRow[i.row.Current]
		i.setLine(0, fmt.Sprintf("Server.Zone.Name:       %v", server.Zone.Name))
		i.setLine(1, fmt.Sprintf("Server.Name:            %v", server.Name))
		i.setLine(2, fmt.Sprintf("Server.Description:     %v", server.Description))
		i.setLine(3, fmt.Sprintf("Server.InterfaceDriver: %v", server.InterfaceDriver))
		i.setLine(4, fmt.Sprintf("Server.ServiceClass:    %v", server.ServiceClass))
		i.setLine(5, fmt.Sprintf("Server.Instance.Status: %v", server.Instance.Status))
		i.setLine(6, fmt.Sprintf("Server.Availability:    %v", server.Availability))
		i.setLine(7, fmt.Sprintf("Server.CreatedAt:       %v", server.CreatedAt))
		i.setLine(8, fmt.Sprintf("Server.ModifiedAt:      %v", server.ModifiedAt))
		i.setLine(9, fmt.Sprintf("Server.Tags:            %v", server.Tags))
		termbox.Flush()
		return
	}

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
		i.setLine(currentRow, server.String(i.unanonymize))
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
		url := i.client.URL(zone, []string{"server"})

		statusCode, respBody, err := i.client.Request("GET", url, nil)
		if err != nil {
			return err
		}

		if statusCode != 200 {
			return errors.New(fmt.Sprintf("Request Method: GET, Request URL: %v, Status Code: %v", url, statusCode))
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
		return fmt.Sprintf("[WARNING] Server.Name %v is already up", s.Name)
	}

	url := i.client.URL(s.Zone.Name, []string{"server", s.ID, "power"})
	statusCode, _, err := i.client.Request("PUT", url, nil)

	if err != nil {
		return fmt.Sprintf("[ERROR] %v", err)
	}

	if statusCode != 200 {
		return fmt.Sprintf("[ERROR] Request Method: PUT, Request URL: %v, Status Code: %v", url, statusCode)
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
