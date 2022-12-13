package models

import (
	"regexp"
	"strings"
)

var link = "\\[.*\\]\\(.*\\)|https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"

var handler *lineHandler

type lineHandler struct {
	regex *regexp.Regexp
}

// TODO move this to configuration
var excludeLinks = []string{
	"http://127.0.0.1:10000",
	"http://127.0.0.1:10001",
}

func GetInstance() *lineHandler {
	regex, _ := regexp.Compile(link)
	if handler == nil {
		handler = &lineHandler{
			regex: regex,
		}
	}
	return handler

}

func (handler *lineHandler) FindAndCheckLinksInLine(path string) []string {
	var linksPaths []string
	var validPaths []string
	linksPaths = append(linksPaths, handler.regex.FindAllString(path, -1)...)
	for _, linkPath := range linksPaths {
		if strings.Contains(linkPath, "](") {
			linkPath = strings.Split(linkPath, "](")[1]
		}
		linkPath = strings.Split(linkPath, ")")[0]
		if !handler.isExcluded(linkPath) {
			validPaths = append(validPaths, linkPath)
		}
	}
	return validPaths
}

func (handler *lineHandler) isExcluded(link string) bool {
	for _, excludedPath := range excludeLinks {
		if strings.HasPrefix(link, excludedPath) {
			return true
		}
	}
	return false
}
