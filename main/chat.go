package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Arindam-langer/OllamaChat/main/ui"
	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
	"github.com/Arindam-langer/OllamaChat/store"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type chatResponseMsg struct {
	userQuery    string
	aiResponse   string
	historyEntry []llms.MessageContent
	err          error
}

// chatDeps holds persistent connections reused across chat messages
type chatDeps struct {
	vectorStore *store.VectorStore
	embedder    *embeddings.EmbedderImpl
	chatLLM     *ollama.LLM
	once        sync.Once
	err         error
}

func (d *chatDeps) init() error {
	d.once.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			d.err = fmt.Errorf("DATABASE_URL is not set in .env")
			return
		}

		embedModel := os.Getenv("EMBED_MODEL")
		if embedModel == "" {
			embedModel = "nomic-embed-text"
		}

		chatModel := os.Getenv("CHAT_MODEL")
		if chatModel == "" {
			chatModel = "qwen2.5:3b"
		}

		ctx := context.Background()
		vs, err := store.Connect(ctx, dbURL)
		if err != nil {
			d.err = fmt.Errorf("failed to connect to vector store: %w", err)
			return
		}
		d.vectorStore = vs

		emb, err := extract.NewOllamaEmbedder(embedModel, "")
		if err != nil {
			d.err = fmt.Errorf("failed to create embedder: %w", err)
			return
		}
		d.embedder = emb

		llm, err := ollama.New(ollama.WithModel(chatModel))
		if err != nil {
			d.err = fmt.Errorf("failed to create chat LLM: %w", err)
			return
		}
		d.chatLLM = llm
	})
	return d.err
}

func doChatCmd(ctx context.Context, deps *chatDeps, query string, history []llms.MessageContent) tea.Cmd {
	return func() tea.Msg {
		if err := deps.init(); err != nil {
			return chatResponseMsg{err: err}
		}

		systemPrompt := os.Getenv("SYSTEM_PROMPT")
		if systemPrompt == "" {
			systemPrompt = `You are a helpful assistant that answers questions based on the provided context documents.
Use the context below to answer the user's question accurately and concisely.
If the context doesn't contain relevant information, say so honestly and try your best to help.`
		}

		topK := 5

		queryVectors, err := deps.embedder.EmbedDocuments(ctx, []string{query})
		if err != nil {
			return chatResponseMsg{err: fmt.Errorf("failed to embed query: %w", err)}
		}
		if len(queryVectors) == 0 {
			return chatResponseMsg{err: fmt.Errorf("no query vector returned")}
		}

		results, err := deps.vectorStore.SearchSimilar(ctx, queryVectors[0], topK)
		if err != nil {
			return chatResponseMsg{err: fmt.Errorf("failed to search similar documents: %w", err)}
		}

		contextStr := ""
		if len(results) > 0 {
			contextStr = strings.Join(results, "\n\n---\n\n")
		}

		userPrompt := query
		if contextStr != "" {
			userPrompt = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextStr, query)
		}

		messages := []llms.MessageContent{
			{
				Role:  llms.ChatMessageTypeSystem,
				Parts: []llms.ContentPart{llms.TextContent{Text: systemPrompt}},
			},
		}
		messages = append(messages, history...)
		messages = append(messages, llms.MessageContent{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: userPrompt}},
		})

		var fullResponse strings.Builder
		_, err = deps.chatLLM.GenerateContent(ctx, messages,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fullResponse.Write(chunk)
				return nil
			}),
		)
		if err != nil {
			return chatResponseMsg{err: fmt.Errorf("failed to generate response: %w", err)}
		}

		historyEntries := []llms.MessageContent{
			{
				Role:  llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{llms.TextContent{Text: query}},
			},
			{
				Role:  llms.ChatMessageTypeAI,
				Parts: []llms.ContentPart{llms.TextContent{Text: fullResponse.String()}},
			},
		}

		return chatResponseMsg{
			userQuery:    query,
			aiResponse:   fullResponse.String(),
			historyEntry: historyEntries,
		}
	}
}

func updateChat(msg tea.Msg, m model) (model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Intercept Enter key message for submission so we don't insert a newline into the textarea
		if key.Matches(msg, ui.Keys.Enter) {
			query := strings.TrimSpace(m.chatInput.Value())
			if query != "" {
				m.chatContent += fmt.Sprintf("\nYou > %s\n", query)
				wrapped := lipgloss.NewStyle().Width(m.chatViewport.Width - 4).Render(m.chatContent)
				m.chatViewport.SetContent(wrapped)
				m.chatViewport.GotoBottom()

				m.chatInput.Reset()
				m.chatLoading = true
				m.chatErr = nil

				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel

				return m, tea.Batch(
					doChatCmd(ctx, m.chatDeps, query, m.chatHistory),
					m.spinner.Tick,
				)
			}
			return m, nil
		}

		// Handle text input typing
		var tiCmd tea.Cmd
		m.chatInput, tiCmd = m.chatInput.Update(msg)
		cmds = append(cmds, tiCmd)

		// Handle viewport navigation keys (pageup/pagedown, etc.)
		var vpCmd tea.Cmd
		m.chatViewport, vpCmd = m.chatViewport.Update(msg)
		cmds = append(cmds, vpCmd)

	case chatResponseMsg:
		m.chatLoading = false
		if msg.err != nil {
			m.chatErr = msg.err
			m.chatContent += fmt.Sprintf("\nError: %v\n", msg.err)
		} else {
			m.chatContent += fmt.Sprintf("\nAI > %s\n\n", msg.aiResponse)
			m.chatHistory = append(m.chatHistory, msg.historyEntry...)
		}
		wrapped := lipgloss.NewStyle().Width(m.chatViewport.Width - 4).Render(m.chatContent)
		m.chatViewport.SetContent(wrapped)
		m.chatViewport.GotoBottom()

	case spinner.TickMsg:
		if m.chatLoading {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			return m, spinCmd
		}
	}

	return m, tea.Batch(cmds...)
}

func viewChat(m model) string {
	var s strings.Builder
	s.WriteString(ui.TitleStyle.Render("  Interactive Chat  "))
	s.WriteString("\n\n")

	// Viewport content
	s.WriteString(m.chatViewport.View())
	s.WriteString("\n\n")

	// Input / loading status line
	if m.chatLoading {
		s.WriteString(m.spinner.View())
		s.WriteString(" Thinking...")
	} else if m.chatErr != nil {
		s.WriteString(ui.CuteHighlight.Render("Error: " + m.chatErr.Error()))
	} else {
		s.WriteString(m.chatInput.View())
	}
	s.WriteString("\n\n")

	return s.String()
}
