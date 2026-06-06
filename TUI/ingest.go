package main

import (
	"strings"

	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func updateIngest(msg tea.Msg, m model) (model, tea.Cmd) {
	return m, nil
}

func viewIngest(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("  Ingest Page  "))
	s.WriteString("\n\n")
	s.WriteString(ui.CuteHighlight.Render("Ingest implementation goes here...\n\n"))
	return s.String()
}
