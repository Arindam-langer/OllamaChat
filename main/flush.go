package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Arindam-langer/OllamaChat/store"
)

func runFlush() {
	fmt.Println("Flushing all data from vector store...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("DATABASE_URL is not set in .ENV")
	}

	ctx := context.Background()

	// Connect to Database
	vectorStore, err := store.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to vector store: %v", err)
	}
	defer vectorStore.Close()

	// Flush the database table
	err = vectorStore.Flush(ctx)
	if err != nil {
		log.Fatalf("Failed to flush vector store: %v", err)
	}

	fmt.Println("Successfully flushed all document embeddings.")
}
