package core

import (
	"log"
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
