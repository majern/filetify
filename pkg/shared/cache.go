package shared

import (
	"github.com/allegro/bigcache/v3"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
	"time"
)

var cacheInstanceLock sync.Once
var cacheInstance *bigcache.BigCache

func getCache() *bigcache.BigCache {
	cacheInstanceLock.Do(func() {
		cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(48 * time.Hour))
		if err != nil {
			log.Fatalf("Unable to add file to cache. Error: %v", err)
		}

		cacheInstance = cache
	})

	return cacheInstance
}

func CacheFile(path string, file FileEntry) error {
	buff := Encode(file)
	return getCache().Set(path, buff)
}

func GetFileFromCache(path string) (*FileEntry, error) {
	if entry, err := getCache().Get(path); err == nil {
		obj := Decode[FileEntry](entry)
		return obj, nil
	} else {
		return nil, err
	}
}

func UpdateFile(path string, operation FileOperation) {
	entry, err := GetFileFromCache(path)
	HandleError(err, true)

	entry.Modified = time.Now()
	entry.Operation = Deleted

	CacheFile(path, *entry)
}

func GetAllFromCache() []*FileEntry {
	var entries []*FileEntry
	cacheIterator := getCache().Iterator()
	for cacheIterator.SetNext() {
		current, err := cacheIterator.Value()

		if err == nil {
			entry := Decode[FileEntry](current.Value())
			entries = append(entries, entry)
			logrus.WithField("entry", entry).Debug("Entry loaded from cache")
		} else {
			logrus.WithError(err).Error("Failed to read entry from cache")
		}
	}

	return entries
}

func GetAllKeysFromCache() []string {
	var keys []string
	cacheIterator := getCache().Iterator()
	for cacheIterator.SetNext() {
		current, err := cacheIterator.Value()

		if err == nil {
			keys = append(keys, current.Key())
			logrus.WithField("key", current.Key()).Debug("Key loaded from cache")
		} else {
			logrus.WithError(err).Error("Failed to read key from cache")
		}
	}

	return keys
}
