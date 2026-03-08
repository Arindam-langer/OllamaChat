/*
extract text\
split them
embed them
store them in db
*/
package extract

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/ledongthuc/pdf"
)

// get a list of file names
func ReadDir(dataSet string) ([]string, error) {
	entries, err := os.ReadDir(dataSet)
	fileNames := []string{}
	if err != nil {
		return nil, fmt.Errorf("Failed to read directory: %w", err)
	}
	for _, file := range entries {
		fileNames = append(fileNames, file.Name())
	}
	log.Printf("Getting files from Directory%s,: %v", dataSet, fileNames)
	return fileNames, nil
}

// get the content of a file
func ReadPdf(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to Open Pdf: %w", err)
	}
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("Failed to Open Pdf: %w", err)
	}
	buf.ReadFrom(b)
	content := buf.String()
	log.Printf("Getting contents of file: %s", filePath)
	return content, nil
}
