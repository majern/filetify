package shared

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ScanFiles(paths []string, serverScan bool, ignoredFilesRegex []string) {
	logrus.Info("Scanning files...")

	if serverScan {
		ScanRecursive(paths[0], paths[0], "/", ignoredFilesRegex, Scanned)
	} else {
		for _, path := range paths {
			dir := filepath.Dir(path)
			parent := strings.ReplaceAll(path, dir, "")
			ScanRecursive(path, path, parent, ignoredFilesRegex, Scanned)
		}
	}

	logrus.Info("File scanning completed.")
}

func ScanRecursive(sourcePath, basePath, relativePath string, ignoredFilesRegex []string, fileOp FileOperation) {
	entries, _ := os.ReadDir(sourcePath)

	for _, entry := range entries {
		entryKey := strings.ReplaceAll(filepath.Join(sourcePath, entry.Name()), basePath, relativePath)
		logger := logrus.WithField("path", entryKey)

		if IsFileIgnored(entryKey, ignoredFilesRegex) {
			continue
		}

		if entry.IsDir() {
			logger.Debug("DIR Scanned")

			ScanRecursive(filepath.Join(sourcePath, entry.Name()), basePath, relativePath, ignoredFilesRegex, fileOp)

		} else {
			logger.Debug("FILE Scanned")
		}

		info, _ := entry.Info()
		fileEntry := NewFileEntry(entryKey, filepath.Join(sourcePath, entry.Name()), entry.IsDir(), info.ModTime().UTC(), fileOp)
		CacheFile(entryKey, fileEntry)
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
