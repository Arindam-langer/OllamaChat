package main

import "github.com/charmbracelet/lipgloss"

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
