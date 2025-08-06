# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- Build: `go build -o ssh-config-manager`
- Run: `./ssh-config-manager` (or `./ssh-config-manager.exe` on Windows)
- Dependency management: `go mod tidy`

## Code Structure

- Main entry: `main.go`
- SSH config handling: `internal/config/ssh_config.go`
- TUI components: `internal/ui/` (model.go, view.go, styles.go)

## Key Libraries

- Bubbletea (TUI framework)
- Lipgloss (UI styling)
- Bubbles (UI components)