package main

import (
	"log"
	"strings"

	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. Handle Global Controls (Force quit, global back to menu, layout resizing)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ui.Keys.ForceQuit):
			return m, tea.Quit

		case key.Matches(msg, ui.Keys.Back):
			m.state = screenMenu
			m.flushing = false
			m.flushSuccess = false
			m.flushError = nil
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.progress.Width = msg.Width - 10
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
	}

	// 2. Delegate state update to active screen helper
	var cmd tea.Cmd
	switch m.state {
	case screenMenu:
		m, cmd = updateMenu(msg, m)
	case screenChat:
		m, cmd = updateChat(msg, m)
	case screenIngest:
		m, cmd = updateIngest(msg, m)
	case screenFlush:
		m, cmd = updateFlush(msg, m)
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	// 1. Render active screen content
	switch m.state {
	case screenMenu:
		s.WriteString(viewMenu(m))
	case screenChat:
		s.WriteString(viewChat(m))
	case screenIngest:
		s.WriteString(viewIngest(m))
	case screenFlush:
		s.WriteString(viewFlush(m))
	}

	// 2. Configure active keys based on current page
	activeKeys := ui.Keys
	if m.state == screenMenu {
		activeKeys.Up.SetEnabled(true)
		activeKeys.Down.SetEnabled(true)
		activeKeys.Enter.SetEnabled(true)
		activeKeys.Back.SetEnabled(false)
		activeKeys.Quit.SetEnabled(true)
		activeKeys.Yes.SetEnabled(false)
		activeKeys.No.SetEnabled(false)
	} else if m.state == screenFlush {
		activeKeys.Up.SetEnabled(false)
		activeKeys.Down.SetEnabled(false)
		activeKeys.Enter.SetEnabled(false)
		activeKeys.Quit.SetEnabled(false)

		if !m.flushing && !m.flushSuccess && m.flushError == nil {
			activeKeys.Yes.SetEnabled(true)
			activeKeys.No.SetEnabled(true)
			activeKeys.Back.SetEnabled(true)
		} else {
			activeKeys.Yes.SetEnabled(false)
			activeKeys.No.SetEnabled(false)
			activeKeys.Back.SetEnabled(true)
		}
	} else {
		activeKeys.Up.SetEnabled(false)
		activeKeys.Down.SetEnabled(false)
		activeKeys.Enter.SetEnabled(false)
		activeKeys.Quit.SetEnabled(false)
		activeKeys.Back.SetEnabled(true)
		activeKeys.Yes.SetEnabled(false)
		activeKeys.No.SetEnabled(false)
	}

	s.WriteString(m.help.View(activeKeys))

	return s.String()
}

func main() {
	godotenv.Load(".env")

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI Error: %v", err)
	}
}
