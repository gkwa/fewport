package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getMarkdownFilePaths(dir string) ([]string, error) {
	var filePaths []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		filePaths = append(filePaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return filterMarkdownFiles(filePaths), nil
}

func filterMarkdownFiles(filePaths []string) []string {
	var markdownFiles []string

	for _, filePath := range filePaths {
		if strings.HasSuffix(filePath, ".md") {
			markdownFiles = append(markdownFiles, filePath)
		}
	}

	return markdownFiles
}

func processFile(path string) ([]ParamInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var newContent strings.Builder
	scanner := bufio.NewScanner(file)
	var remainingParams []ParamInfo

	for scanner.Scan() {
		line := scanner.Text()
		newLine, modified, params := cleanLine(line, path)
		newContent.WriteString(newLine + "\n")

		if modified {
			fmt.Printf("Modified URL in %s\n", path)
		}

		remainingParams = append(remainingParams, params...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, []byte(newContent.String()), 0o644); err != nil {
		return nil, err
	}

	return remainingParams, nil
}
