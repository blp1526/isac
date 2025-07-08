# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building and Testing
- `make build` - Builds the project (runs linting, tests, and readme generation)
- `make test` - Runs lint checks and tests with coverage
- `make lint` - Runs code quality checks (`go vet` and `gofmt`)
- `go test -v --cover ./...` - Run tests directly with coverage

### Code Quality
- `go vet ./...` - Static analysis tool (replaces gometalinter)
- `gofmt -l .` - Code formatting check
- The project uses standard Go tools for code quality (migrated from gometalinter)

### Other Commands
- `make readme` - Generates README from template
- `make clean` - Removes dist/ directory
- `make release` - Uses goreleaser for releases

## Architecture Overview

### Core Components
- **main.go**: Entry point using urfave/cli v3 Command-based architecture with context.Context
- **lib/cmd/cmd.go**: CLI application setup using urfave/cli/v3 library with NewCommand() function
- **lib/isac.go**: Main TUI application logic with termbox-based interface
- **lib/api/client.go**: HTTP client for SAKURA Cloud API interactions
- **lib/config/config.go**: Configuration management for API credentials and zones

### TUI Framework
- Uses `termbox-go` for terminal UI rendering
- Main event loop handles keyboard input and screen updates
- Key components:
  - **lib/row/row.go**: Row management and cursor positioning
  - **lib/state/state.go**: Application state management (help, detail views)
  - **lib/keybinding/keybinding.go**: Keyboard shortcut definitions

### Data Models
- **lib/resource/server/server.go**: SAKURA Cloud server resource representation
- Server data is fetched from multiple zones and displayed in a filterable table

### Application Flow
1. Configuration loaded from `$HOME/.usacloud/default/config.json`
2. CLI Command created with urfave/cli v3 architecture
3. API client initialized with access tokens
4. TUI starts with server list from specified zones
5. User can navigate, filter, and interact with servers
6. Real-time server operations (power on) via API calls

### Key Features
- Multi-zone server listing and filtering
- Interactive server management (power on/off)
- Real-time status updates
- Keyboard-driven navigation
- Server detail views

## Development Notes

### Modern Go Practices
- **Go Modules**: Uses go.mod instead of legacy dep/Gopkg.toml
- **urfave/cli v3**: Modern Command-based CLI architecture with context.Context support
- **Standard Tooling**: Uses `go vet` and `gofmt` instead of deprecated gometalinter

### Configuration
- Uses `.usacloud/default/config.json` for API credentials
- Supports multiple zones via CLI flags or config
- Config file created via `isac init` command

### Testing
- Unit tests exist for core components (client, server, row, state)
- Test files follow `*_test.go` naming convention
- Uses standard Go testing framework

### CLI Framework Migration Notes
- Migrated from urfave/cli v1 → v2 → v3
- v3 uses `&cli.Command{}` instead of `cli.NewApp()`
- Action functions now receive `context.Context` as first parameter
- Authors field expects string format: `"Name <email>"`