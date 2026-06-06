package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
	"github.com/Arindam-langer/OllamaChat/store"
	"github.com/Arindam-langer/OllamaChat/main/ui"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ingestFinishedMsg struct {
	result string
	err    error
}

// performIngestion runs the full ingest pipeline and returns a summary string
func performIngestion(ctx context.Context) (string, error) {
	embedModel := os.Getenv("EMBED_MODEL")
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "", fmt.Errorf("DATABASE_URL is not set in .env")
	}

	dataset := os.Getenv("DATASET")
	if dataset == "" {
		return "", fmt.Errorf("DATASET path is not set in .env")
	}

	vectorStore, err := store.Connect(ctx, dbURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to vector store: %w", err)
	}
	defer vectorStore.Close()

	if err := vectorStore.InitDB(ctx); err != nil {
		return "", fmt.Errorf("failed to initialize database schema: %w", err)
	}

	embedder, err := extract.NewOllamaEmbedder(embedModel, "")
	if err != nil {
		return "", fmt.Errorf("failed to create embedder: %w", err)
	}

	fileNames, err := extract.ReadDir(dataset)
	if err != nil {
		return "", fmt.Errorf("failed to read dataset directory: %w", err)
	}

	var totalFiles int
	var totalChunks int

	for _, file := range fileNames {
		if filepath.Ext(file) != ".pdf" {
			continue
		}
		filePath := filepath.Join(dataset, file)

		content, err := extract.ReadPdf(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read PDF %s: %w", file, err)
		}

		chunks, err := extract.ChunkText(content, 800, 120)
		if err != nil {
			return "", fmt.Errorf("failed to chunk text %s: %w", file, err)
		}

		embeddings, err := extract.EmbedChunks(ctx, embedder, chunks)
		if err != nil {
			return "", fmt.Errorf("failed to embed %s: %w", file, err)
		}

		err = vectorStore.InsertEmbeddings(ctx, embeddings)
		if err != nil {
			return "", fmt.Errorf("failed to store embeddings for %s: %w", file, err)
		}

		totalFiles++
		totalChunks += len(embeddings)
	}

	return fmt.Sprintf("Successfully processed %d files and stored %d embeddings.", totalFiles, totalChunks), nil
}

func doIngestCmd() tea.Cmd {
	return func() tea.Msg {
		res, err := performIngestion(context.Background())
		return ingestFinishedMsg{result: res, err: err}
	}
}

func updateIngest(msg tea.Msg, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.ingesting && !m.ingestDone && m.ingestErr == nil {
			switch {
			case key.Matches(msg, ui.Keys.Yes):
				m.ingesting = true
				m.ingestDone = false
				m.ingestErr = nil
				m.ingestResult = ""
				return m, tea.Batch(
					doIngestCmd(),
					m.spinner.Tick,
				)
			case key.Matches(msg, ui.Keys.No):
				m.state = screenMenu
				return m, nil
			}
		}

	case ingestFinishedMsg:
		m.ingesting = false
		if msg.err != nil {
			m.ingestErr = msg.err
			return m, nil
		}
		m.ingestDone = true
		m.ingestResult = msg.result
		return m, nil

	case spinner.TickMsg:
		if m.ingesting {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func viewIngest(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("  Dataset Ingestion  "))
	s.WriteString("\n\n")

	if !m.ingesting && !m.ingestDone && m.ingestErr == nil {
		s.WriteString("Ready to ingest PDF documents from the dataset directory?\n")
		s.WriteString("This will read, chunk, embed, and store them in the database.\n\n")
	} else if m.ingesting {
		s.WriteString(m.spinner.View())
		s.WriteString(" Ingesting dataset... Please wait.\n\n")
	} else if m.ingestDone {
		s.WriteString("Dataset ingestion complete!\n\n")
		s.WriteString(ui.UnselectedStyle.Render(m.ingestResult))
		s.WriteString("\n\n")
	} else if m.ingestErr != nil {
		s.WriteString("Error during ingestion:\n")
		s.WriteString(ui.CuteHighlight.Render(m.ingestErr.Error()))
		s.WriteString("\n\n")
	}
	return s.String()
}
