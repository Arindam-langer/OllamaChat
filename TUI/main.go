package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89B4FA")) // Soft bright blue

	// Selected menu item
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A6E3A1")) // Soft green

	// Unselected menu item
	unselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BAC2DE")) // Muted gray-blue

	// Secondary text / highlights
	cuteHighlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5C2E7")) // Soft pink
)

// 1. Define Keybindings using bubbles/key
type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	Back      key.Binding
	Quit      key.Binding
	ForceQuit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view (not used here but required by interface).
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Back, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc/backspace", "back to menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "force quit"),
	),
}

type screenState int

const (
	screenMenu screenState = iota
	screenChat
	screenIngest
	screenFlush
)

// 2. Add help.Model to the Model struct
type model struct {
	options []string
	cursor  int
	state   screenState
	help    help.Model
}

func initialModel() model {
	h := help.New()
	// Style keys with cute pink, and descriptions with soft slate blue
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5C2E7"))
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#585B70"))

	return model{
		options: []string{"run chat", "ingest", "Flush DB", "exit"},
		cursor:  0,
		state:   screenMenu,
		help:    h,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// 3. Update key matches in the Update Method
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
	}
	return m, nil
}

// 4. Render the dynamic Help View
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

	// Dynamic help: disable Up/Down/Enter bindings if not on Menu state
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
