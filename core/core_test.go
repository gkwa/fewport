package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanGoogleURLs(t *testing.T) {
	// Create test markdown files
	testDir := filepath.Join("testdata", "markdown")
	err := os.MkdirAll(testDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll("testdata")

	file1 := filepath.Join(testDir, "file1.md")
	file1Content := `
# File 1

This is a test file with a Google URL: https://www.google.com/search?q=test&client=chrome&ie=UTF-8&og=og
`
	err = os.WriteFile(file1, []byte(file1Content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	file2 := filepath.Join(testDir, "file2.md")
	file2Content := `
# File 2

This is another test file with a Google URL: https://www.google.com/search?q=example&client=firefox&ie=UTF-8&og=og
`
	err = os.WriteFile(file2, []byte(file2Content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Run CleanGoogleURLs
	err = CleanGoogleURLs(testDir)
	if err != nil {
		t.Fatalf("CleanGoogleURLs failed: %v", err)
	}

	// Check if the URLs were cleaned
	cleanedFile1, _ := os.ReadFile(file1)
	expectedFile1 := `
# File 1

This is a test file with a Google URL: https://www.google.com/search?q=test
`
	if string(cleanedFile1) != expectedFile1 {
		t.Errorf("File 1 content mismatch. Expected:\n%s\nGot:\n%s", expectedFile1, string(cleanedFile1))
	}

	cleanedFile2, _ := os.ReadFile(file2)
	expectedFile2 := `
# File 2

This is another test file with a Google URL: https://www.google.com/search?q=example
`
	if string(cleanedFile2) != expectedFile2 {
		t.Errorf("File 2 content mismatch. Expected:\n%s\nGot:\n%s", expectedFile2, string(cleanedFile2))
	}
}
