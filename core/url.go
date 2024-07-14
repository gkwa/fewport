package core

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

var excludePatterns = []string{
	"google.com/maps/",
}

var excludeParams = []string{
	"q",
	"tbm",
}

func cleanLine(line, filePath string) (string, bool, []ParamInfo) {
	rxStrict := xurls.Strict()
	urls := rxStrict.FindAllString(line, -1)

	modifiedLine := line
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
			modifiedLine = strings.ReplaceAll(modifiedLine, u, cleanedURL)

			for _, param := range params {
				remainingParams = append(remainingParams, ParamInfo{
					Param:     param,
					FilePath:  filePath,
					SourceURL: u,
				})
			}
		}
	}

	if modifiedLine != line {
		modified = true
	}

	return modifiedLine, modified, remainingParams
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
		"fbs",
		"gs_lcrp",
		"gs_lp",
		"gs_lcp",
		"gs_ssp",
		"ictx",
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
		// "udm", // no! don't remove this one, udm=2 means its an image search, eg https://www.google.com/search?udm=2&q=poison+ivy
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
