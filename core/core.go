package core

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"mvdan.cc/xurls/v2"
)

//go:embed report.md
var templateFS embed.FS

type ParamInfo struct {
	Param     string
	FilePath  string
	SourceURL string
	SearchURL string
}

type ReportData struct {
	ParamsToRemove []string
	ParamInfoList  []ParamInfo
}

var excludePatterns = []string{
	"google.com/maps/",
}

var excludeParams = []string{
	"q",
	"tbm",
}

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

func updateParamInfoMap(paramInfoMap map[string]ParamInfo, param ParamInfo) {
	if contains(excludeParams, param.Param) {
		return
	}

	if paramInfo, ok := paramInfoMap[param.Param]; ok {
		paramInfo.FilePath = param.FilePath
		paramInfoMap[param.Param] = paramInfo
	} else {
		searchQuery := fmt.Sprintf("what is purpose of google url parameter \"%s\"", param.Param)
		searchURL := fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(searchQuery))
		paramInfoMap[param.Param] = ParamInfo{
			Param:     param.Param,
			FilePath:  param.FilePath,
			SourceURL: param.SourceURL,
			SearchURL: searchURL,
		}
	}
}

func generateReport(paramInfoMap map[string]ParamInfo) error {
	tmpl, err := template.ParseFS(templateFS, "report.md")
	if err != nil {
		return err
	}

	var paramsToRemove []string
	for param := range paramInfoMap {
		paramsToRemove = append(paramsToRemove, param)
	}
	sort.Strings(paramsToRemove)

	var paramInfoList []ParamInfo
	for _, paramInfo := range paramInfoMap {
		paramInfoList = append(paramInfoList, paramInfo)
	}

	reportData := ReportData{
		ParamsToRemove: paramsToRemove,
		ParamInfoList:  paramInfoList,
	}

	reportFile, err := os.Create("remaining_params_report.md")
	if err != nil {
		return err
	}
	defer reportFile.Close()

	err = tmpl.Execute(reportFile, reportData)
	if err != nil {
		return err
	}

	log.Printf("Remaining parameters report generated: remaining_params_report.md")
	return nil
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

func cleanLine(line, filePath string) (string, bool, []ParamInfo) {
	rxStrict := xurls.Strict()
	urls := rxStrict.FindAllString(line, -1)

	modified := false
	var remainingParams []ParamInfo

	for _, u := range urls {
		if isExcludedURL(u) {
			continue
		}

		if strings.Contains(strings.ToLower(u), "google.com") {
			cleanedURL, params, err := cleanGoogleURL(u)
			if err != nil {
				log.Printf("Failed to clean URL %s: %v", u, err)
				continue
			}
			line = strings.ReplaceAll(line, u, cleanedURL)
			modified = true

			for _, param := range params {
				remainingParams = append(remainingParams, ParamInfo{
					Param:     param,
					FilePath:  filePath,
					SourceURL: u,
				})
			}
		}
	}

	return line, modified, remainingParams
}

func isExcludedURL(urlStr string) bool {
	for _, pattern := range excludePatterns {
		if strings.Contains(strings.ToLower(urlStr), pattern) {
			return true
		}
	}
	return false
}

func cleanGoogleURL(urlStr string) (string, []string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", nil, err
	}

	q := u.Query()
	paramsToRemove := []string{
		"bih",
		"biw",
		"client",
		"dpr",
		"ei",
		"gs_lcrp",
		"gs_lp",
		"ie",
		"oq",
		"prmd",
		"sa",
		"sca_esv",
		"sca_upv",
		"sclient",
		"source",
		"sourceid",
		"sqi",
		"sxsrf",
		"uact",
		"udm",
		"uds",
		"ved",
	}

	var remainingParams []string

	for param := range q {
		if contains(excludeParams, param) || !contains(paramsToRemove, param) {
			remainingParams = append(remainingParams, param)
			continue
		}
		q.Del(param)
	}

	u.RawQuery = q.Encode()
	return u.String(), remainingParams, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
