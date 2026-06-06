package main

import (
	"fmt"
	"strings"

	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func updateMenu(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ui.Keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, ui.Keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, ui.Keys.Down):
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case key.Matches(msg, ui.Keys.Enter):
			switch m.cursor {
			case 0:
				m.state = screenChat
			case 1:
				m.state = screenIngest
				m.ingesting = false
				m.ingestDone = false
				m.ingestErr = nil
				m.ingestResult = ""
			case 2:
				m.state = screenFlush
				m.flushing = false
				m.flushSuccess = false
				m.flushError = nil
				m.progressVal = 0.0
			case 3: // exit
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func viewMenu(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("  Ollama Chat  "))
	s.WriteString("\n\n")
	for i, option := range m.options {
		if m.cursor == i {
			s.WriteString(ui.SelectedStyle.Render(fmt.Sprintf("%s", option)))
			s.WriteString("\n")
		} else {
			s.WriteString(ui.UnselectedStyle.Render(fmt.Sprintf("  %s", option)))
			s.WriteString("\n")
		}
	}
	s.WriteString("\n")
	return s.String()
}
