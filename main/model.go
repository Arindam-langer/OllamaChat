package main

import (
	"github.com/Arindam-langer/OllamaChat/main/ui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmc/langchaingo/llms"
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

	// Chat state fields
	chatDeps     *chatDeps
	chatInput    textarea.Model
	chatViewport viewport.Model
	chatLoading  bool
	chatErr      error
	chatHistory  []llms.MessageContent
	chatContent  string
}

func initialModel() model {
	h := help.New()
	h.Styles.ShortKey = ui.CuteHighlight.Copy()
	h.Styles.ShortDesc = ui.DimStyle.Copy()

	prog := progress.New(progress.WithDefaultGradient())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#89B4FA"))

	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Focus()
	ta.CharLimit = 500
	ta.SetWidth(60)
	ta.SetHeight(3)

	vp := viewport.New(80, 15)
	vp.SetContent("Ask anything about your ingested documents.")

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
		chatDeps:     &chatDeps{},
		chatInput:    ta,
		chatViewport: vp,
		chatLoading:  false,
		chatErr:      nil,
		chatHistory:  nil,
		chatContent:  "Ask anything about your ingested documents.\n\n",
	}
}
