package core

import (
   "os"
   "path/filepath"
   "testing"
)

func TestCleanGoogleURLs(t *testing.T) {
   testCases := []struct {
   	name           string
   	inputContent   string
   	expectedOutput string
   }{
   	{
   		name: "File 1",
   		inputContent: `
# File 1
This is a test file with a Google URL: https://www.google.com/search?q=test&client=chrome&ie=UTF-8&og=stuff
`,
   		expectedOutput: `
# File 1
This is a test file with a Google URL: https://www.google.com/search?og=stuff&q=test
`,
   	},
   	{
   		name: "File 2",
   		inputContent: `
# File 2
This is another test file with a Google URL: https://www.google.com/search?q=example&client=firefox&ie=UTF-8&og=og
`,
   		expectedOutput: `
# File 2
This is another test file with a Google URL: https://www.google.com/search?og=og&q=example
`,
   	},
   }

   for _, tc := range testCases {
   	t.Run(tc.name, func(t *testing.T) {
   		testDir := filepath.Join("testdata", "markdown")
   		err := os.MkdirAll(testDir, 0o755)
   		if err != nil {
   			t.Fatalf("Failed to create test directory: %v", err)
   		}
   		defer os.RemoveAll("testdata")

   		file := filepath.Join(testDir, tc.name+".md")
   		err = os.WriteFile(file, []byte(tc.inputContent), 0o644)
   		if err != nil {
   			t.Fatalf("Failed to write file: %v", err)
   		}

   		err = CleanGoogleURLs(testDir)
   		if err != nil {
   			t.Fatalf("CleanGoogleURLs failed: %v", err)
   		}

   		cleanedContent, _ := os.ReadFile(file)
   		if string(cleanedContent) != tc.expectedOutput {
   			t.Errorf("Content mismatch. Expected:\n%s\nGot:\n%s", tc.expectedOutput, string(cleanedContent))
   		}
   	})
   }
}
