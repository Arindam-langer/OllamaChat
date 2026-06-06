package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.ForceQuit):
			return m, tea.Quit

		case key.Matches(msg, keys.Quit):
			if m.state == screenMenu {
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Back):
			m.state = screenMenu
			return m, nil
		}

		// Handle keys based on the active screen
		switch m.state {
		case screenMenu:
			switch {
			case key.Matches(msg, keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}
			case key.Matches(msg, keys.Down):
				if m.cursor < len(m.options)-1 {
					m.cursor++
				}
			case key.Matches(msg, keys.Enter):
				switch m.cursor {
				case 0:
					m.state = screenChat
				case 1:
					m.state = screenIngest
				case 2:
					m.state = screenFlush
				case 3: // exit
					return m, tea.Quit
				}
			}
		}

	case tea.WindowSizeMsg:
		// Let the help model know the terminal width so it can wrap/truncate
		m.help.Width = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	switch m.state {
	case screenMenu:
		s.WriteString(titleStyle.Render("  Ollama Chat  "))
		s.WriteString("\n\n")
		for i, option := range m.options {
			if m.cursor == i {
				s.WriteString(selectedStyle.Render(fmt.Sprintf("%s", option)))
				s.WriteString("\n")
			} else {
				s.WriteString(unselectedStyle.Render(fmt.Sprintf("  %s", option)))
				s.WriteString("\n")
			}
		}
		s.WriteString("\n")

	case screenChat:
		s.WriteString(titleStyle.Render("Run Chat Page  "))
		s.WriteString("\n\n")
		s.WriteString(cuteHighlight.Render("Chat implementation goes here...\n\n"))

	case screenIngest:
		s.WriteString(titleStyle.Render("  Ingest Page  "))
		s.WriteString("\n\n")
		s.WriteString(cuteHighlight.Render("Ingest implementation goes here...\n\n"))

	case screenFlush:
		s.WriteString(titleStyle.Render("  Flush DB Page  "))
		s.WriteString("\n\n")
		s.WriteString(cuteHighlight.Render("Flush implementation goes here...\n\n"))
	}

	// Dynamic help: disable Up/Down/Enter/Quit bindings if not on Menu state
	activeKeys := keys
	if m.state != screenMenu {
		activeKeys.Up.SetEnabled(false)
		activeKeys.Down.SetEnabled(false)
		activeKeys.Enter.SetEnabled(false)
		activeKeys.Quit.SetEnabled(false) // q shouldn't quit from subpages, only esc
		activeKeys.Back.SetEnabled(true)
	} else {
		activeKeys.Back.SetEnabled(false) // Esc has no action on main menu
	}

	// Render the help view
	s.WriteString(m.help.View(activeKeys))

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI Error: %v", err)
	}
}
