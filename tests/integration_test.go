package tests

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/shiranr/linkcheck/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFilesWithError(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	logFilePath := setLogFile()
	basepath := filepath.Dir(b)
	configPath := basepath + "/resources/linkcheck.json"
	utils.LoadConfiguration(configPath)
	readmeFiles := utils.ExtractMarkdownFiles()
	logFileContent := readLogFile(logFilePath)
	res := models.GetFilesProcessorInstance().Process(readmeFiles)
	assert.ErrorContains(t, res, "ERROR: 3 links check failed, please check the logs")
	assert.Contains(t, logFileContent, "Went through 2 files")
	assert.Contains(t, logFileContent, "Line 21 link nla.go status 400")
	assert.Contains(t, logFileContent, "Line 33 link source-control/merge-strategies.md status 400")
	assert.Contains(t, logFileContent, "Line 33 link resources/templates/CONTRIBUTING.md status 400")
	assert.NotContains(t, logFileContent, "http://bla.com/")
	assert.NotContains(t, logFileContent, "http://test.com/")

}

func readLogFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func setLogFile() string {
	logFilePath := "resources/log.txt"
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).
			Fatal("Failed to open log file.")
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	return logFilePath
}
