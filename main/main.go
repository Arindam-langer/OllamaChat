package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
	"github.com/joho/godotenv"
)

const DATASET string = "Test_dataset"

func main() {
	fmt.Println("Running the program")
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Println("No .ENV file found, using default environment variables")
	}

	embedModel := os.Getenv("EMBED_MODEL")
	ctx := context.Background()
	embedder, err := extract.NewOllamaEmbedder(embedModel, "")
	if err != nil {
		log.Fatalln("Failed to create embedder:", err)
	}

	fileNames, err := extract.ReadDir(DATASET)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range fileNames {
		filePath := filepath.Join(DATASET, file)
		content, err := extract.ReadPdf(filePath)
		if err != nil {
			panic(err)
		}
		chunks, err := extract.ChunkText(content, 800, 120)
		if err != nil {
			panic(err)
		}
		fmt.Println("Total chunks:", len(chunks))

		// embed the chunks
		embeddings, err := extract.EmbedChunks(ctx, embedder, chunks)
		if err != nil {
			log.Fatalf("Failed to embed %s: %v", file, err)
		}
		fmt.Printf("File: %s → %d embeddings (dim: %d)\n", file, len(embeddings), len(embeddings[0].Vector))
	}
}
