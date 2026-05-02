package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
	"github.com/Arindam-langer/OllamaChat/store"
	"github.com/joho/godotenv"
)

func main() {
	// Parse basic CLI commands
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Load environment variables
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Println("No .ENV file found, using default environment variables")
	}

	switch command {
	case "ingest":
		runIngest()
	case "chat":
		runChat()
	case "show":
		runShow()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Ollama Chat - Local RAG Pipeline")
	fmt.Println("Usage:")
	fmt.Println("  go run ./main ingest   - Reads PDFs from dataset, embeds them, and stores in Postgres")
	fmt.Println("  go run ./main chat     - Start the chat UI (Not yet implemented)")
	fmt.Println("  go run ./main show     - Show all documents")
}

func runIngest() {
	fmt.Println("Starting Ingestion Process...")

	embedModel := os.Getenv("EMBED_MODEL")
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("DATABASE_URL is not set in .ENV")
	}

	dataset := os.Getenv("DATASET")
	if dataset == "" {
		dataset = "Test_dataset" // Default fallback
	}

	ctx := context.Background()

	//Connect to Database
	fmt.Println("Connecting to Postgres...")
	vectorStore, err := store.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to vector store: %v", err)
	}
	defer vectorStore.Close()

	if err := vectorStore.InitDB(ctx); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	//Setup Embedder
	embedder, err := extract.NewOllamaEmbedder(embedModel, "")
	if err != nil {
		log.Fatalln("Failed to create embedder:", err)
	}

	//Process Files
	fileNames, err := extract.ReadDir(dataset)
	if err != nil {
		log.Fatalln("Failed to read dataset directory:", err)
	}

	for _, file := range fileNames {
		filePath := filepath.Join(dataset, file)
		fmt.Printf("Processing %s...\n", file)

		content, err := extract.ReadPdf(filePath)
		if err != nil {
			log.Printf("Failed to read PDF %s: %v\n", file, err)
			continue
		}

		chunks, err := extract.ChunkText(content, 800, 120)
		if err != nil {
			log.Printf("Failed to chunk text %s: %v\n", file, err)
			continue
		}
		fmt.Println("  Total chunks:", len(chunks))

		// embed the chunks
		embeddings, err := extract.EmbedChunks(ctx, embedder, chunks)
		if err != nil {
			log.Printf("Failed to embed %s: %v\n", file, err)
			continue
		}

		// store embeddings in postgres
		err = vectorStore.InsertEmbeddings(ctx, embeddings)
		if err != nil {
			log.Printf("Failed to store embeddings for %s: %v\n", file, err)
			continue
		}

		fmt.Printf("  Successfully ingested %s -> %d embeddings (dim: %d)\n", file, len(embeddings), len(embeddings[0].Vector))
	}

	fmt.Println("Ingestion complete.")
}

func runChat() {
	fmt.Println("Chat interface is not yet implemented.")
	// will be coming in future
}
func runShow() {
	dataset := os.Getenv("DATASET")
	dir, err := os.ReadDir(dataset)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range dir {
		fmt.Printf("Name: %s | IsDir: %v\n", entry.Name(), entry.IsDir())
	}
}
