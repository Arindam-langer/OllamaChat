package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
	"github.com/Arindam-langer/OllamaChat/store"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func runChat() {

	//remove fallbacks once the chat loops are deterministic
	chatModel := os.Getenv("CHAT_MODEL")
	if chatModel == "" {
		chatModel = "qwen2.5:3b"
	}

	embedModel := os.Getenv("EMBED_MODEL")
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}
	systemPrompt := os.Getenv("SYSTEM_PROMPT")
	if systemPrompt == "" {
		systemPrompt = `You are a helpful assistant that answers questions based on the provided context documents.
Use the context below to answer the user's question accurately and concisely.
If the context doesn't contain relevant information, say so honestly and try your best to help.`
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("DATABASE_URL is not set in .ENV")
	}

	topK := 5

	ctx := context.Background()

	fmt.Println("Connecting to vector store...")
	vectorStore, err := store.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to vector store: %v", err)
	}
	defer vectorStore.Close()

	if err := vectorStore.InitDB(ctx); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Setup embedder for query embedding (same model used during ingestion)
	embedder, err := extract.NewOllamaEmbedder(embedModel, "")
	if err != nil {
		log.Fatalf("Failed to create embedder: %v", err)
	}

	// Setup chat LLM
	llm, err := ollama.New(ollama.WithModel(chatModel))
	if err != nil {
		log.Fatalf("Failed to create chat LLM: %v", err)
	}

	fmt.Printf("\n OllamaChat (model: %s) — type your question, or 'quit' to exit\n\n", chatModel)

	scanner := bufio.NewScanner(os.Stdin)
	var history []llms.MessageContent

	for {
		fmt.Print("You > ")
		if !scanner.Scan() {
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}
		if query == "quit" || query == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		queryVectors, err := embedder.EmbedDocuments(ctx, []string{query})
		if err != nil {
			log.Printf("Failed to embed query: %v\n", err)
			continue
		}
		if len(queryVectors) == 0 {
			log.Println("No query vector returned")
			continue
		}

		// Retrieve similar chunks from pgvector
		results, err := vectorStore.SearchSimilar(ctx, queryVectors[0], topK)
		if err != nil {
			log.Printf("Failed to search similar documents: %v\n", err)
			continue
		}

		contextStr := ""
		if len(results) > 0 {
			contextStr = strings.Join(results, "\n\n---\n\n")
		}

		userPrompt := query
		if contextStr != "" {
			userPrompt = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextStr, query)
		}

		// 5. Assemble messages: system + conversation history + current query
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

		fmt.Print("\nAI > ")
		var fullResponse strings.Builder

		_, err = llm.GenerateContent(ctx, messages,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				fullResponse.Write(chunk)
				return nil
			}),
		)
		if err != nil {
			log.Printf("\nFailed to generate response: %v\n", err)
			continue
		}
		fmt.Print("\n\n")

		history = append(history,
			llms.MessageContent{
				Role:  llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{llms.TextContent{Text: query}},
			},
			llms.MessageContent{
				Role:  llms.ChatMessageTypeAI,
				Parts: []llms.ContentPart{llms.TextContent{Text: fullResponse.String()}},
			},
		)
	}
}
