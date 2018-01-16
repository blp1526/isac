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
	"github.com/blp1526/isac/lib/state"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

type Isac struct {
	filter         string
	client         *api.Client
	config         *config.Config
	showCurrentRow bool
	unanonymize    bool
	row            *row.Row
	// FIXME: wasteful
	serverByCurrentRow map[int]server.Server
	servers            []server.Server
	state              *state.State
	reverseSort        bool
	zones              []string
	message            string
}

func New(configPath string, unanonymize bool, zones string) (i *Isac, err error) {
	config, err := config.New(configPath)
	if err != nil {
		return i, err
	}

	state := state.New()

	client := api.NewClient(config.AccessToken, config.AccessTokenSecret)
	row := row.New()

	zs := []string{config.Zone}
	if zones != "" {
		zs = strings.Split(zones, ",")
	}

	i = &Isac{
		client:         client,
		config:         config,
		showCurrentRow: true,
		unanonymize:    unanonymize,
		row:            row,
		state:          state,
		zones:          zs,
		reverseSort:    false,
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
				i.currentServerUp()
			case termbox.KeyCtrlSlash:
				i.state.Toggle("help")
				i.draw("")
			case termbox.KeyBackspace2, termbox.KeyCtrlB, termbox.KeyCtrlH:
				i.removeRuneFromFilter()
			case termbox.KeyCtrlS:
				i.reverseSort = !i.reverseSort
				i.draw("")
			case termbox.KeyEnter:
				i.state.Toggle("detail")
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

		if i.showCurrentRow && i.row.Current == y {
			fgColor = termbox.ColorBlack
			bgColor = termbox.ColorYellow
		}

		termbox.SetCell(x, y, r, fgColor, bgColor)
		x += runewidth.RuneWidth(r)
	}
}

func (i *Isac) draw(message string) {
	termbox.Clear(coldef, coldef)

	if i.state.Current == "help" {
		i.showCurrentRow = false
		lines := []string{
			"Quick reference for isac keybindings:",
			"",
			"<ESC>, <C-c>            exit",
			"<Arrow Up>, <C-p>       move current row up",
			"<Arrow Down>, <C-n>     move current row down",
			"<C-u>                   power on current row's server",
			"<C-r>                   refresh rows",
			"<BackSpace, C-b>, <C-h> delete a filter character",
			"<C-s>                   sort rows",
			"<C-/>                   show help",
			"<Enter>                 show current row's detail",
		}

		for index, line := range lines {
			i.setLine(index, line)
		}

		termbox.Flush()
		return
	}

	if i.state.Current == "detail" {
		i.showCurrentRow = false
		server := i.serverByCurrentRow[i.row.Current]

		lines := []string{
			fmt.Sprintf("Server.Zone.Name:       %v", server.Zone.Name),
			fmt.Sprintf("Server.Name:            %v", server.Name),
			fmt.Sprintf("Server.Description:     %v", server.Description),
			fmt.Sprintf("Server.InterfaceDriver: %v", server.InterfaceDriver),
			fmt.Sprintf("Server.ServiceClass:    %v", server.ServiceClass),
			fmt.Sprintf("Server.Instance.Status: %v", server.Instance.Status),
			fmt.Sprintf("Server.Availability:    %v", server.Availability),
			fmt.Sprintf("Server.CreatedAt:       %v", server.CreatedAt),
			fmt.Sprintf("Server.ModifiedAt:      %v", server.ModifiedAt),
			fmt.Sprintf("Server.Tags:            %v", server.Tags),
		}

		for index, line := range lines {
			i.setLine(index, line)
		}

		termbox.Flush()
		return
	}

	i.showCurrentRow = true

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
		}

		return servers[x].Instance.Status < servers[y].Instance.Status
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

func (i *Isac) currentServerUp() {
	s := i.serverByCurrentRow[i.row.Current]

	if s.ID == "" {
		i.draw("[ERROR] Current row has no Server")
		return
	}

	if s.Instance.Status == "up" {
		i.draw(fmt.Sprintf("[WARNING] Server.Name %v is already up", s.Name))
		return
	}

	go func() {
		url := i.client.URL(s.Zone.Name, []string{"server", s.ID, "power"})
		statusCode, _, err := i.client.Request("PUT", url, nil)

		if err != nil {
			i.draw(fmt.Sprintf("[ERROR] %v", err))
			return
		}

		if statusCode != 200 {
			i.draw(fmt.Sprintf("[ERROR] Request Method: PUT, Server.Name: %v, Status Code: %v", s.Name, statusCode))
			return
		}

		i.draw(fmt.Sprintf("Server.Name %v is booting, wait few seconds, and refresh", s.Name))
		return
	}()
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
