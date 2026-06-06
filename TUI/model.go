package main

import (
	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
)

type screenState int

const (
	screenMenu screenState = iota
	screenChat
	screenIngest
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
}

func initialModel() model {
	h := help.New()
	// Style keys and descriptions using the theme styles
	h.Styles.ShortKey = ui.CuteHighlight.Copy()
	h.Styles.ShortDesc = ui.DimStyle.Copy()

	prog := progress.New(progress.WithDefaultGradient())

	return model{
		options:      []string{"run chat", "ingest", "Flush DB", "exit"},
		cursor:       0,
		state:        screenMenu,
		help:         h,
		progress:     prog,
		progressVal:  0.0,
		flushing:     false,
		dbDone:       false,
		flushSuccess: false,
		flushError:   nil,
	}
}
