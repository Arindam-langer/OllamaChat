package main

import (
	"fmt"
	"log"
	"os"

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
