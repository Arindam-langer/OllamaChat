package main

import (
	"fmt"
	"log"
	"path/filepath"

	extract "github.com/Arindam-langer/OllamaChat/preprocessing"
)

const DATASET string = "Test_dataset"

func main() {
	fmt.Println("Running the program")
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
	}
}
