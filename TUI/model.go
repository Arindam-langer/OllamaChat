package main

import (
	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type screenState int

const (
	screenMenu screenState = iota
	screenChat
	screenIngest
	screenShow
	screenFlush
)

type model struct {
	options      []string
	cursor       int
	state        screenState
	help         help.Model
	progress     progress.Model
	progressVal  float64
	flushing     bool
	dbDone       bool
	flushSuccess bool
	flushError   error

	// Ingest state fields
	spinner      spinner.Model
	ingesting    bool
	ingestDone   bool
	ingestErr    error
	ingestResult string

	// Show state fields
	showFiles  []string
	showCursor int
	showErr    error
}

func initialModel() model {
	h := help.New()
	h.Styles.ShortKey = ui.CuteHighlight.Copy()
	h.Styles.ShortDesc = ui.DimStyle.Copy()

	prog := progress.New(progress.WithDefaultGradient())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#89B4FA"))

	return model{
		options:      []string{"Run Chat", "Ingest", "Show Files", "Flush DB", "Exit"},
		cursor:       0,
		state:        screenMenu,
		help:         h,
		progress:     prog,
		progressVal:  0.0,
		flushing:     false,
		dbDone:       false,
		flushSuccess: false,
		flushError:   nil,
		spinner:      s,
		ingesting:    false,
		ingestDone:   false,
		ingestErr:    nil,
		ingestResult: "",
	}
}
