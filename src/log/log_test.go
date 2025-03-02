package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetCacheDir(t *testing.T) {
	// Test when AWS_ECR_CACHE_DIR is not set
	os.Unsetenv("AWS_ECR_CACHE_DIR")
	expected := "~/.sm"
	actual := GetCacheDir()
	assert.Equal(t, expected, actual, "they should be equal")

	// Test when AWS_ECR_CACHE_DIR is set
	expected = "/tmp/cache"
	os.Setenv("AWS_ECR_CACHE_DIR", expected)
	actual = GetCacheDir()
	assert.Equal(t, expected, actual, "they should be equal")
}

func TestLogrusConfig(t *testing.T) {
	// Set up a temporary directory for testing
	tempDir := t.TempDir()
	os.Setenv("AWS_ECR_CACHE_DIR", tempDir)

	// Call the function
	LogrusConfig()

	// Verify the log file is created
	logFilePath := filepath.Join(tempDir, "log", "sm-login.log")
	_, err := os.Stat(logFilePath)
	assert.False(t, os.IsNotExist(err), "log file should exist")

	// Verify logrus configuration
	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel(), "log level should be Debug")
	file, err := os.OpenFile(logFilePath, os.O_RDONLY, 0644)
	assert.NoError(t, err, "should be able to open log file")
	defer file.Close()
}

func TestLogrusConfigWithError(t *testing.T) {
	// Set up an invalid directory to force an error
	invalidDir := string([]byte{0})
	os.Setenv("AWS_ECR_CACHE_DIR", invalidDir)

	// Call the function
	LogrusConfig()

	// Verify logrus configuration falls back to default
	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel(), "log level should be Debug")
}