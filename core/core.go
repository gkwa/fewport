package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CleanGoogleURLs(dir string) error {
	var failedFiles []string
	paramInfoMap := make(map[string]ParamInfo)

	filePaths, err := getMarkdownFilePaths(dir)
	if err != nil {
		return err
	}

	for _, filePath := range filePaths {
		params, err := processFile(filePath)
		if err != nil {
			log.Printf("Failed to process file %s: %v", filePath, err)
			failedFiles = append(failedFiles, filePath)
			continue
		}

		for _, param := range params {
			updateParamInfoMap(paramInfoMap, param)
		}
	}

	if len(failedFiles) > 0 {
		log.Printf("Failed to process the following files:")
		for _, file := range failedFiles {
			log.Printf("- %s", file)
		}
	}

	if len(paramInfoMap) > 0 {
		err := generateReport(paramInfoMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func CleanGoogleURLsInFile(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		newLine, _, _ := cleanLine(line, "")
		fmt.Fprintln(w, newLine)
	}
	return scanner.Err()
}

func ProcessPathsFromStdin(ctx context.Context, transform func(io.Reader, io.Writer) error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		path := scanner.Text()
		err := ProcessFile(ctx, path, transform)
		if err != nil {
			log.Printf("Failed to process file %s: %v", path, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

func ProcessFile(ctx context.Context, path string, transform func(io.Reader, io.Writer) error) error {
	originalContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	var processedContent strings.Builder
	err = transform(bytes.NewReader(originalContent), &processedContent)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	if bytes.Equal(originalContent, []byte(processedContent.String())) {
		return nil
	}

	err = os.WriteFile(path, []byte(processedContent.String()), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write processed content to file: %w", err)
	}

	return nil
}
