package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)
func GetCacheDir() string {
	if cacheDir := os.Getenv("AWS_ECR_CACHE_DIR"); cacheDir != "" {
		return cacheDir
	}
	return "~/.sm"
}
func LogrusConfig()() {
	logdir, err := homedir.Expand(GetCacheDir() + "/log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to find directory: %v", err)
		logdir = os.TempDir()
	}
	// Clean the path to replace with OS-specific separators
	logdir = filepath.Clean(logdir)
	err = os.MkdirAll(logdir, os.ModeDir|0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to create directory: %v", err)
		logdir = os.TempDir()
	}
	file, err := os.OpenFile(filepath.Join(logdir, "sm-login.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(file)
}