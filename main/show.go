package main

import (
	"fmt"
	"log"
	"os"
)

func runShow() {
	dataset := os.Getenv("DATASET")
	dir, err := os.ReadDir(dataset)
	if err != nil {
		log.Fatal(err)
	}
	for i, entry := range dir {
		fmt.Printf("%d. Name: %s\n", i+1, entry.Name())
	}
}
