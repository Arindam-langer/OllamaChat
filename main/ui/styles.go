package ui

import "github.com/charmbracelet/lipgloss"

var (
	// TitleStyle is the style for page and app headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89B4FA")) // Soft bright blue

	// SelectedStyle is the style for selected options
	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A6E3A1")) // Soft green

	// UnselectedStyle is the style for unselected options
	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BAC2DE")) // Muted gray-blue

	// CuteHighlight is the style for secondary highlights / success states
	CuteHighlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5C2E7")) // Soft pink

	// DimStyle is the style for secondary descriptions / instructions
	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#585B70")) // Soft slate gray
)
