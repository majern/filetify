package shared

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func ScanFiles(paths []string) {
	logrus.Info("Scanning files...")
	ignoredFilesRegex := GetConfiguration().IgnoredFiles

	for _, path := range paths {
		//ScanSinglePath(path)
		scanRecursive(path, ignoredFilesRegex)
	}

	logrus.Info("File scanning completed.")
}

func scanRecursive(path string, ignoredFilesRegex []string) {
	entries, _ := os.ReadDir(path)

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		logger := logrus.WithField("path", entryPath)

		if IsFileIgnored(entryPath, ignoredFilesRegex) {
			continue
		}

		if entry.IsDir() {
			logger.Debug("DIR Scanned")

			scanRecursive(entryPath, ignoredFilesRegex)

		} else {
			logger.Debug("FILE Scanned")
		}

		entry := NewFileEntry(entryPath, entry.IsDir(), time.Now(), Scanned)
		CacheFile(entryPath, *entry)
	}
}

func IsFileIgnored(path string, ignoredFilesRegex []string) bool {
	fileIgnored := false

	for _, regex := range ignoredFilesRegex {
		match, _ := regexp.MatchString(regex, path)
		if match == true && fileIgnored == false {
			fileIgnored = true
		}
	}

	return fileIgnored
}
