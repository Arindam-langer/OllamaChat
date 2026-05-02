/*
extract text\
split them
embed them
store them in db
*/
package preprocessing

import (
	"fmt"
	"log"
	"os"

	"github.com/ledongthuc/pdf"
	"github.com/tmc/langchaingo/textsplitter"
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
	content := ""
	defer f.Close()

	totalPages := r.NumPage()

	for i := range totalPages {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("Failed to get page number %d: %w", i, err)
		}
		content += text //keep it inefficient for now but we need a refactor later here.

	}
	log.Printf("Getting contents of file: %s", filePath)
	return content, nil
}

// chunk text
func ChunkText(text string, size int, overlap int) ([]string, error) {

	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(size),
		textsplitter.WithChunkOverlap(overlap),
	)

	chunks, err := splitter.SplitText(text)
	if err != nil {
		return nil, fmt.Errorf("Failed to chunk text %w", err)
	}

	return chunks, nil
}
