package core

import (
	"embed"
	"log"
	"os"
	"sort"
	"text/template"
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
