package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

type screenState int

const (
	screenMenu screenState = iota
	screenChat
	screenIngest
	screenFlush
)

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
