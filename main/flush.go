package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Arindam-langer/OllamaChat/main/ui"
	"github.com/Arindam-langer/OllamaChat/store"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type progressTickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressTickMsg(t)
	})
}

type flushFinishedMsg struct {
	err error
}

func doFlushCmd(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			return flushFinishedMsg{err: fmt.Errorf("DATABASE_URL is not set in .env")}
		}
		vectorStore, err := store.Connect(ctx, dbURL)
		if err != nil {
			return flushFinishedMsg{err: err}
		}
		defer vectorStore.Close()
		err = vectorStore.Flush(ctx)
		return flushFinishedMsg{err: err}
	}
}

func updateFlush(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.flushing && !m.flushSuccess && m.flushError == nil {
			switch {
			case key.Matches(msg, ui.Keys.Yes):
				m.flushing = true
				m.dbDone = false
				m.flushSuccess = false
				m.flushError = nil
				m.progressVal = 0.0
				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel
				return m, tea.Batch(
					doFlushCmd(ctx),
					tickCmd(),
				)
			case key.Matches(msg, ui.Keys.No):
				m.state = screenMenu
				return m, nil
			}
		}

	case progressTickMsg:
		if !m.flushing {
			return m, nil
		}
		m.progressVal += 0.1
		if m.progressVal >= 0.9 && !m.dbDone {
			m.progressVal = 0.9
		}
		if m.progressVal >= 1.0 {
			m.progressVal = 1.0
			m.flushing = false
			m.flushSuccess = true
			return m, nil
		}
		return m, tickCmd()

	case flushFinishedMsg:
		if msg.err != nil {
			m.flushError = msg.err
			m.flushing = false
			return m, nil
		}
		m.dbDone = true
		m.progressVal = 1.0
		m.flushing = false
		m.flushSuccess = true
		return m, nil
	}
	return m, nil
}

func viewFlush(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("  Flush Database  "))
	s.WriteString("\n\n")

	if !m.flushing && !m.flushSuccess && m.flushError == nil {
		s.WriteString("Are you sure you want to delete all document embeddings?\n")
		s.WriteString("This action cannot be undone.\n\n")
	} else if m.flushing {
		s.WriteString("Flushing database... Please wait.\n\n")
		s.WriteString(m.progress.ViewAs(m.progressVal))
		s.WriteString("\n\n")
	} else if m.flushSuccess {
		s.WriteString("Database successfully flushed!\n\n")
		s.WriteString(ui.UnselectedStyle.Render("All documents and embeddings have been deleted."))
		s.WriteString("\n\n")
	} else if m.flushError != nil {
		s.WriteString("Error flushing database:\n")
		s.WriteString(ui.CuteHighlight.Render(m.flushError.Error()))
		s.WriteString("\n\n")
	}
	return s.String()
}
