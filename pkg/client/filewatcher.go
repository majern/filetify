package client

import (
	"github.com/fsnotify/fsnotify"
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var watcher *fsnotify.Watcher

func WatchFilesInPath(dir string) {
	logrus.Info("Initializing file watcher")
	watcher, _ = fsnotify.NewWatcher()

	defer func(watcher *fsnotify.Watcher) {
		shared.HandleError(watcher.Close(), false)
		logrus.Warning("File watcher stopped")
	}(watcher)

	logrus.Info("Adding directories to file watcher...")
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
	logrus.WithField("file", event.Name).Info("File DELETED") //Send deleted due to https://github.com/fsnotify/fsnotify/issues/529

	fileEntry, err := shared.GetFileFromCache(event.Name)
	shared.HandleError(err, true)

	if fileEntry.IsDir {
		relatedKeys := shared.GetAllKeysFromCache()
		for _, key := range relatedKeys {
			if strings.Contains(key, fileEntry.Path) {
				shared.UpdateFile(key, shared.Deleted)
				logrus.WithField("file", key).Info("File DELETED")
			}
		}
	}

	shared.UpdateFile(fileEntry.Path, shared.Deleted)
}

func handleCreated(event *fsnotify.Event, ignoredFilesRegex []string) {
	logrus.WithField("file", event.Name).Info("File CREATED") //Send new

	if isDir(event.Name) {
		shared.ScanFiles([]string{event.Name}, ignoredFilesRegex)
		shared.HandleErrorWithMsg(filepath.Walk(event.Name, watchDir), true, "Failed to add directory to file watcher")
		shared.CacheFile(event.Name, shared.NewFileEntry(event.Name, true, time.Now(), shared.New))
	} else {
		shared.CacheFile(event.Name, shared.NewFileEntry(event.Name, false, time.Now(), shared.New))
	}
}

func handleUpdated(event *fsnotify.Event) {
	logrus.WithField("file", event.Name).Info("File UPDATED")

	shared.UpdateFile(event.Name, shared.Modified)
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

func isDir(path string) bool {
	stat, err := os.Stat(path)
	shared.HandleErrorWithMsg(err, true, "Failed to read file info")
	return stat.IsDir()
}
