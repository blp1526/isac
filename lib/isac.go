package isac

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/blp1526/isac/lib/api"
	"github.com/blp1526/isac/lib/config"
	"github.com/blp1526/isac/lib/keybinding"
	"github.com/blp1526/isac/lib/resource/server"
	"github.com/blp1526/isac/lib/row"
	"github.com/blp1526/isac/lib/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the state of the application following Bubble Tea's MVU pattern
type Model struct {
	filter         string
	client         *api.Client
	config         *config.Config
	showCurrentRow bool
	row            *row.Row
	// FIXME: wasteful
	serverByCurrentRow map[int]server.Server
	servers            []server.Server
	state              *state.State
	reverseSort        bool
	zones              []string
	message            string
	ready              bool
	width              int
	height             int
}

// serverUpdateMsg is used to update servers in the background
type serverUpdateMsg struct {
	servers []server.Server
	err     error
}

// statusMsg is used to show status messages
type statusMsg string

// Init is called when the program starts
func (m Model) Init() tea.Cmd {
	return m.loadServers()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "ctrl+p":
			m.currentRowUp()
			return m, nil

		case "down", "ctrl+n":
			m.currentRowDown()
			return m, nil

		case "ctrl+r":
			return m, m.loadServers()

		case "ctrl+u":
			return m, m.currentServerUp()

		case "ctrl+/":
			m.state.Toggle("help")
			return m, nil

		case "ctrl+a":
			m.removeRuneAllFromFilter()
			return m, nil

		case "backspace", "ctrl+b", "ctrl+h":
			m.removeRuneFromFilter()
			return m, nil

		case "ctrl+s":
			m.reverseSort = !m.reverseSort
			return m, nil

		case "enter":
			m.state.Toggle("detail")
			return m, nil

		default:
			// Handle text input for filtering
			if len(msg.Runes) > 0 {
				m.addRuneToFilter(msg.Runes[0])
			}
			return m, nil
		}

	case serverUpdateMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("[ERROR] %v", msg.err)
		} else {
			m.servers = msg.servers
			m.message = "Servers have been refreshed"
		}
		return m, nil

	case statusMsg:
		m.message = string(msg)
		return m, nil
	}

	return m, nil
}

// View renders the interface
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var content string

	// Show help screen
	if m.state.Current == "help" {
		lines := []string{
			"Quick reference for isac keybindings:",
			"",
		}
		for _, k := range keybinding.Keybindings() {
			lines = append(lines, fmt.Sprintf("%s %s", k.Keys, k.Desc))
		}
		content = strings.Join(lines, "\n")
		return content
	}

	// Show detail screen
	if m.state.Current == "detail" {
		if server, exists := m.serverByCurrentRow[m.row.Current]; exists {
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
			content = strings.Join(lines, "\n")
			return content
		}
	}

	// Main view - server list
	m.showCurrentRow = true

	// Filter servers
	var filteredServers []server.Server
	for _, s := range m.servers {
		if strings.Contains(s.Name, m.filter) {
			filteredServers = append(filteredServers, s)
		}
	}

	// Sort servers
	sort.Slice(filteredServers, func(x, y int) bool {
		if m.reverseSort {
			return filteredServers[x].Instance.Status > filteredServers[y].Instance.Status
		}
		return filteredServers[x].Instance.Status < filteredServers[y].Instance.Status
	})

	// Update current row bounds
	if m.row.Current == 0 {
		m.row.Current = m.row.HeadersSize()
	}
	m.row.MovableBottom = len(filteredServers) + m.row.HeadersSize() - 1
	if m.row.Current > m.row.MovableBottom {
		m.row.Current = m.row.HeadersSize()
	}

	// Build headers
	headers := m.row.Headers(m.message, strings.Join(m.zones, ", "), len(filteredServers), m.currentNo(), m.filter)

	// Build server list
	var lines []string
	lines = append(lines, headers...)

	m.serverByCurrentRow = map[int]server.Server{}

	// Style for current row
	currentRowStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("11")).
		Foreground(lipgloss.Color("0"))

	for index, server := range filteredServers {
		currentRow := index + m.row.HeadersSize()
		line := server.String()

		// Highlight current row
		if m.showCurrentRow && m.row.Current == currentRow {
			line = currentRowStyle.Render(line)
		}

		lines = append(lines, line)
		m.serverByCurrentRow[currentRow] = server
	}

	content = strings.Join(lines, "\n")

	// Ensure we don't exceed terminal height
	if len(lines) > m.height-1 {
		lines = lines[:m.height-1]
		content = strings.Join(lines, "\n")
	}

	return content
}

// New initializes a new Model
func New(configPath string, zones string) (*Model, error) {
	config, err := config.New(configPath)
	if err != nil {
		return nil, err
	}

	state := state.New()
	client := api.NewClient(config.AccessToken, config.AccessTokenSecret)
	row := row.New()

	zs := []string{config.Zone}
	if zones != "" {
		zs = strings.Split(zones, ",")
	}

	model := &Model{
		client:         client,
		config:         config,
		showCurrentRow: true,
		row:            row,
		state:          state,
		zones:          zs,
		reverseSort:    false,
	}

	return model, nil
}

// Run starts the Bubble Tea program
func (m *Model) Run() error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Helper methods

func (m *Model) currentRowUp() {
	if m.row.Current > m.row.HeadersSize() {
		m.row.Current--
	}
}

func (m *Model) currentRowDown() {
	if m.row.Current < m.row.MovableBottom {
		m.row.Current++
	}
}

func (m *Model) currentNo() int {
	return m.row.Current + 1 - m.row.HeadersSize()
}

func (m *Model) addRuneToFilter(r rune) {
	m.filter = m.filter + string(r)
}

func (m *Model) removeRuneFromFilter() {
	r := []rune(m.filter)
	if len(r) > 0 {
		m.filter = string(r[:(len(r) - 1)])
	}
}

func (m *Model) removeRuneAllFromFilter() {
	m.filter = ""
}

// loadServers returns a command to load servers from the API
func (m Model) loadServers() tea.Cmd {
	return func() tea.Msg {
		var servers []server.Server

		for _, zone := range m.zones {
			url := m.client.URL(zone, []string{"server"})

			statusCode, respBody, err := m.client.Request("GET", url, nil)
			if err != nil {
				return serverUpdateMsg{servers: nil, err: err}
			}

			if statusCode != 200 {
				return serverUpdateMsg{
					servers: nil,
					err:     fmt.Errorf("Request Method: GET, Request URL: %v, Status Code: %v", url, statusCode),
				}
			}

			sc := server.NewCollection(zone)
			err = json.Unmarshal(respBody, sc)
			if err != nil {
				return serverUpdateMsg{servers: nil, err: err}
			}

			for _, s := range sc.Servers {
				servers = append(servers, s)
			}
		}

		return serverUpdateMsg{servers: servers, err: nil}
	}
}

// currentServerUp returns a command to power on the current server
func (m Model) currentServerUp() tea.Cmd {
	s := m.serverByCurrentRow[m.row.Current]

	if s.ID == "" {
		return func() tea.Msg {
			return statusMsg("[ERROR] Current row has no Server")
		}
	}

	if s.Instance.Status == "up" {
		return func() tea.Msg {
			return statusMsg(fmt.Sprintf("[WARNING] Server.Name %v is already up", s.Name))
		}
	}

	return func() tea.Msg {
		url := m.client.URL(s.Zone.Name, []string{"server", s.ID, "power"})
		statusCode, _, err := m.client.Request("PUT", url, nil)

		if err != nil {
			return statusMsg(fmt.Sprintf("[ERROR] %v", err))
		}

		if statusCode != 200 {
			return statusMsg(fmt.Sprintf("[ERROR] Request Method: PUT, Server.Name: %v, Status Code: %v", s.Name, statusCode))
		}

		return statusMsg(fmt.Sprintf("Server.Name %v is booting, wait few seconds, and refresh", s.Name))
	}
}
