package main

import (
	"fmt"
	"log"
	"path/filepath"

	extract "github.com/Arindam-langer/OllamaChat/Extract"
)

const DATASET string = "dataset"

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
		fmt.Println(content)
	}
}
