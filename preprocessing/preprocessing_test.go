package preprocessing

import (
	"context"
	"os"
	"slices"
	"testing"

	"github.com/joho/godotenv"
)

const dataSet string = "../Test_dataset"

func TestReadDir(t *testing.T) {
	got, _ := ReadDir(dataSet)
	want := []string{"testing.pdf", "testing_2.pdf"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v wanted %v", got, want)
	}
}

func TestEmbedChunks(t *testing.T) {
	_ = godotenv.Load("../.env")
	embedModel := os.Getenv("EMBED_MODEL")
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}

	embedder, err := NewOllamaEmbedder(embedModel, "")
	if err != nil {
		t.Skipf("Skipping test, could not initialize ollama embedder (is ollama running?): %v", err)
	}

	chunks := []string{
		"This is the first test chunk.",
		"This is the second test chunk for embedding.",
	}

	ctx := context.Background()
	embeddings, err := EmbedChunks(ctx, embedder, chunks)
	if err != nil {
		t.Fatalf("EmbedChunks failed: %v", err)
	}

	if len(embeddings) != len(chunks) {
		t.Errorf("Expected %d embeddings, got %d", len(chunks), len(embeddings))
	}

	for i, emb := range embeddings {
		if emb.Text != chunks[i] {
			t.Errorf("Expected chunk text %q, got %q", chunks[i], emb.Text)
		}
		if len(emb.Vector) == 0 {
			t.Errorf("Expected vector for chunk %d to be non-empty", i)
		}
	}
}
