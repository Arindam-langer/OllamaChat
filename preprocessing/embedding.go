package preprocessing

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Embedding holds a text chunk and its vector representation
type Embedding struct {
	Text   string
	Vector []float32
}

func NewOllamaEmbedder(model string, serverURL string) (*embeddings.EmbedderImpl, error) {
	opts := []ollama.Option{
		ollama.WithModel(model),
	}
	// no need for url since we are using localhost
	if serverURL != "" {
		opts = append(opts, ollama.WithServerURL(serverURL))
	}

	llm, err := ollama.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	return embedder, nil
}

func EmbedChunks(ctx context.Context, embedder *embeddings.EmbedderImpl, chunks []string) ([]Embedding, error) {
	if len(chunks) == 0 {
		return nil, nil
	}

	vectors, err := embedder.EmbedDocuments(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("failed to embed chunks: %w", err)
	}

	if len(vectors) != len(chunks) {
		return nil, fmt.Errorf("mismatch: got %d vectors for %d chunks", len(vectors), len(chunks))
	}

	result := make([]Embedding, len(chunks))
	for i, chunk := range chunks {
		result[i] = Embedding{
			Text:   chunk,
			Vector: vectors[i],
		}
	}

	log.Printf("Embedded %d chunks (vector dim: %d)", len(result), len(result[0].Vector))
	return result, nil
}
