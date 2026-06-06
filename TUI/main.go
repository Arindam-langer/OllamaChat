package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// 1. The Model
type model struct {
	options []string
	cursor  int
	state   screenState
}

// screen state to track the page
type screenState int

const (
	screenMenu screenState = iota
	screenChat
	screenIngest
	screenFlush
)

func initialModel() model {
	return model{
		options: []string{"run chat", "ingest", "Flush DB", "exit"},
		cursor:  0,
		state:   screenMenu,
	}
}

// 2. Init Method
func (m model) Init() tea.Cmd {
	return nil
}

// 3. Update Method
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Only quit on 'q' if we are on the menu
			if m.state == screenMenu {
				return m, tea.Quit
			}

		case "esc":
			// Go back to the menu from any screen
			m.state = screenMenu
			return m, nil
		}

		// Handle keys based on the active screen
		switch m.state {
		case screenMenu:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.options)-1 {
					m.cursor++
				}
			case "enter":
				// Transition to the selected screen
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
	}
	return m, nil
}

// 4. View Method
func (m model) View() string {
	var s strings.Builder
	switch m.state {
	case screenMenu:
		s.WriteString("=== Ollama Chat ===\n\n")
		for i, option := range m.options {
			if m.cursor == i {
				s.WriteString(fmt.Sprintf("> %s\n", option))
			} else {
				s.WriteString(fmt.Sprintf("  %s\n", option))
			}
		}
		s.WriteString("\nUse j/k to move • Enter to select • q to quit")
	case screenChat:
		s.WriteString("=== Run Chat Page ===\n\n")
		s.WriteString("Chat implementation goes here...\n\n")
		s.WriteString("Press Esc to return to menu")
	case screenIngest:
		s.WriteString("=== Ingest Page ===\n\n")
		s.WriteString("Ingest implementation goes here...\n\n")
		s.WriteString("Press Esc to return to menu")
	case screenFlush:
		s.WriteString("=== Flush DB Page ===\n\n")
		s.WriteString("Flush implementation goes here...\n\n")
		s.WriteString("Press Esc to return to menu")
	}
	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI Error: %v", err)
	}
}
