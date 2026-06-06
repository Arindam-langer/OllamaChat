package main

import (
	"strings"

	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func updateChat(msg tea.Msg, m model) (model, tea.Cmd) {
	return m, nil
}

func viewChat(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("Run Chat Page  "))
	s.WriteString("\n\n")
	s.WriteString(ui.CuteHighlight.Render("Chat implementation goes here...\n\n"))
	return s.String()
}
