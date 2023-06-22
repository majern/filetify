package client

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var watcher *fsnotify.Watcher
var basePaths map[string]string

func WatchFilesInPath(dir string) {
	logrus.Info("Initializing file watcher")
	watcher, _ = fsnotify.NewWatcher()

	defer func(watcher *fsnotify.Watcher) {
		shared.HandleError(watcher.Close(), false)
		logrus.Warning("File watcher stopped")
	}(watcher)

	logrus.Info("Adding directories to file watcher...")
	basePaths = make(map[string]string)
	for _, path := range GetConfiguration().Paths {
		dir := filepath.Dir(path)
		parent := strings.ReplaceAll(path, dir, "")
		basePaths[path] = parent
	}

	if err := filepath.Walk(dir, watchDir); err != nil {
		logrus.WithError(err).Error("Failed to add directories to file watcher")
		panic(err)
	}

	logrus.Info("Directories added to file watcher")

	//
	done := make(chan bool)
	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if !strings.Contains(event.Name, "~") {
					handleOperation(event)
				}
				// watch for errors
			case err := <-watcher.Errors:
				logrus.WithError(err).Error("Failed to read event from file watcher")
				panic(err)
			}
		}
	}()

	<-done
}

func handleOperation(event fsnotify.Event) {
	ignoredFilesRegex := GetConfiguration().IgnoredFiles
	deleted := fsnotify.Remove | fsnotify.Rename

	if !shared.IsFileIgnored(event.Name, ignoredFilesRegex) {
		switch event.Op {
		case fsnotify.Write:
			handleUpdated(&event)
		case fsnotify.Create:
			handleCreated(&event, ignoredFilesRegex)
		case fsnotify.Remove, fsnotify.Rename, deleted:
			handleDeleted(&event)

			//When directory is deleted, then get all keys from cache, find all that contains directory path and set as deleted
			//After synchronization is completed, remove deleted keys from cache
		default:
			logrus.WithFields(logrus.Fields{
				"file": event.Name,
				"op":   event.Op.String(),
			}).Debug("Unsupported event")
		}
	}
}

func handleDeleted(event *fsnotify.Event) {
	relative, _, _ := getPaths(event.Name)
	fileEntry, err := shared.GetFileFromCache(relative)
	shared.HandleError(err, true)

	logrus.WithField("path", relative).Info(fileOrDir(fileEntry.IsDir) + "DELETED") //Send deleted due to https://github.com/fsnotify/fsnotify/issues/529

	if fileEntry.IsDir {
		relatedKeys := shared.GetAllKeysFromCache()
		for _, key := range relatedKeys {
			if strings.Contains(key, relative) {
				shared.UpdateFile(key, shared.Deleted, time.Now().UTC())
				logrus.WithField("key", key).Info("KEY marked as DELETED")
			}
		}
	}

	shared.UpdateFile(relative, shared.Deleted, time.Now().UTC())
}

func handleCreated(event *fsnotify.Event, ignoredFilesRegex []string) {
	relative, parent, base := getPaths(event.Name)
	fileInfo, _ := os.Stat(event.Name)

	logrus.WithField("path", relative).Info(fileOrDir(fileInfo.IsDir()) + "CREATED") //Send new

	if fileInfo.IsDir() {
		shared.ScanRecursive(event.Name, base, parent, ignoredFilesRegex, shared.New)
		shared.HandleErrorWithMsg(filepath.Walk(event.Name, watchDir), true, "Failed to add directory to file watcher")
		shared.CacheFile(relative, shared.NewFileEntry(relative, event.Name, true, fileInfo.ModTime().UTC(), shared.New))
	} else {
		shared.CacheFile(relative, shared.NewFileEntry(relative, event.Name, false, fileInfo.ModTime().UTC(), shared.New))
	}
}

func handleUpdated(event *fsnotify.Event) {
	relative, _, _ := getPaths(event.Name)
	fileInfo, _ := os.Stat(event.Name)
	logrus.WithField("path", relative).Info(fileOrDir(fileInfo.IsDir()) + "UPDATED")

	shared.UpdateFile(relative, shared.Modified, fileInfo.ModTime().UTC())
}

// watchDir gets run as a walk func, searching for directories to add watchers to.
func watchDir(path string, fi os.FileInfo, err error) error {
	ignoredFilesRegex := GetConfiguration().IgnoredFiles

	if fi.Mode().IsDir() && !shared.IsFileIgnored(path, ignoredFilesRegex) {
		var err = watcher.Add(path)

		if err != nil {
			logrus.WithError(err).WithField("dir", path).Error("An error occured whle adding directory to watcher")
			panic(err)
		} else {
			logrus.Debug("Path added to watcher", path)
		}

		return err
	}

	return nil
}

func getPaths(path string) (relative, parent, base string) {
	for k, v := range basePaths {
		if strings.Contains(path, k) {
			return strings.ReplaceAll(path, k, v), v, k
		}
	}

	shared.HandleErrorWithMsg(errors.New(fmt.Sprintf("Cannot find relative path for physical path='%v'", path)), true, "Error while getting relative path")
	return path, "", ""
}

func isDir(path string) bool {
	stat, err := os.Stat(path)
	shared.HandleErrorWithMsg(err, true, "Failed to read file info")
	return stat.IsDir()
}

func fileOrDir(isDir bool) string {
	if isDir {
		return "DIR "
	}

	return "FILE "
}
