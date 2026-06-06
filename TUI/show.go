package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Arindam-langer/OllamaChat/TUI/ui"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type showLoadedMsg struct {
	files []string
	err   error
}

type openFileMsg struct {
	err error
}

func doShowCmd() tea.Cmd {
	return func() tea.Msg {
		dataset := os.Getenv("DATASET")
		if dataset == "" {
			return showLoadedMsg{err: fmt.Errorf("DATASET path is not set in .env")}
		}
		dir, err := os.ReadDir(dataset)
		if err != nil {
			return showLoadedMsg{err: fmt.Errorf("failed to read dataset directory: %w", err)}
		}
		var files []string
		for _, entry := range dir {
			if !entry.IsDir() {
				files = append(files, entry.Name())
			}
		}
		return showLoadedMsg{files: files}
	}
}

func openFileCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("xdg-open", filePath)
		err := cmd.Start()
		return openFileMsg{err: err}
	}
}

func updateShow(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(m.showFiles) > 0 {
			switch {
			case key.Matches(msg, ui.Keys.Up):
				if m.showCursor > 0 {
					m.showCursor--
				}
			case key.Matches(msg, ui.Keys.Down):
				if m.showCursor < len(m.showFiles)-1 {
					m.showCursor++
				}
			case key.Matches(msg, ui.Keys.Enter):
				dataset := os.Getenv("DATASET")
				selected := m.showFiles[m.showCursor]
				fullPath := filepath.Join(dataset, selected)
				return m, openFileCmd(fullPath)
			}
		}

	case showLoadedMsg:
		if msg.err != nil {
			m.showErr = msg.err
			return m, nil
		}
		m.showFiles = msg.files
		m.showCursor = 0
		return m, nil

	case openFileMsg:
		if msg.err != nil {
			m.showErr = msg.err
		}
		return m, nil
	}
	return m, nil
}

func viewShow(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("  Dataset Files  "))
	s.WriteString("\n\n")

	if m.showFiles == nil && m.showErr == nil {
		s.WriteString("Loading files...\n\n")
	} else if m.showErr != nil {
		s.WriteString("Error reading dataset:\n")
		s.WriteString(ui.CuteHighlight.Render(m.showErr.Error()))
		s.WriteString("\n\n")
	} else if len(m.showFiles) == 0 {
		s.WriteString(ui.UnselectedStyle.Render("No files found in dataset directory."))
		s.WriteString("\n\n")
	} else {
		for i, file := range m.showFiles {
			if i == m.showCursor {
				s.WriteString(ui.SelectedStyle.Render(fmt.Sprintf("  > %s", file)))
			} else {
				s.WriteString(ui.UnselectedStyle.Render(fmt.Sprintf("    %s", file)))
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(ui.DimStyle.Render(fmt.Sprintf("  %d files  |  enter to open", len(m.showFiles))))
		s.WriteString("\n\n")
	}
	return s.String()
}
