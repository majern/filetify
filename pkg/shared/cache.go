package shared

import (
	"github.com/allegro/bigcache/v3"
	proto_v1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"sync"
	"time"
)

const cacheFile = "filetify.cache"

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

func CacheFile(path string, file *FileEntry) error {
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

func UpdateFile(path string, operation FileOperation, modTime time.Time) {
	entry, err := GetFileFromCache(path)
	HandleError(err, true)

	entry.Modified = modTime.UTC()
	entry.Operation = operation

	HandleError(CacheFile(path, entry), true)
}

func GetAllFromCache() []*FileEntry {
	var entries []*FileEntry
	cacheIterator := getCache().Iterator()
	for cacheIterator.SetNext() {
		current, err := cacheIterator.Value()

		if err == nil {
			entry := Decode[FileEntry](current.Value())
			entries = append(entries, entry)
			logrus.WithFields(logrus.Fields{"key": current.Key(), "path": entry.Path, "operation": entry.Operation, "modified": entry.Modified.Format("2006-01-02 15:04:05")}).Debug("Entry loaded from cache")
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

func DumpToFile() {
	logrus.Info("Writing cache dump into file...")
	var entries [][]byte
	cacheIterator := getCache().Iterator()
	for cacheIterator.SetNext() {
		current, err := cacheIterator.Value()
		HandleError(err, true)
		entries = append(entries, current.Value())
	}
	cacheDump := &proto_v1.CacheDump{Files: entries}
	out, err := proto.Marshal(cacheDump)
	HandleErrorWithMsg(err, true, "Failed to serialize cache dump")

	err = os.WriteFile(cacheFile, out, 0644)
	HandleErrorWithMsg(err, true, "Failed to write cache dump into file")

	logrus.WithField("file", cacheFile).Info("Cache dump written into file")
}

func readCacheDumpFromFile() []*FileEntry {
	logrus.WithField("file", cacheFile).Info("Reading cache dump from file...")
	in, err := os.ReadFile(cacheFile)
	HandleErrorWithMsg(err, true, "Failed to read cache dump from file")

	cacheDump := &proto_v1.CacheDump{}
	err = proto.Unmarshal(in, cacheDump)
	HandleErrorWithMsg(err, true, "Failed to deserialize cache dump")

	var entries []*FileEntry
	for _, file := range cacheDump.Files {
		entry := Decode[FileEntry](file)
		entries = append(entries, entry)
	}

	logrus.Info("Cache dump loaded from file")

	return entries
}

// UpdateCacheFromDump is basically used for mark fils and dirs as removed in the cache.
func UpdateCacheFromDump() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithError(r.(error)).Warning("Failed to read cache dump file. Update cache from dump file skipped.")
		}
	}()
	dumpEntries := readCacheDumpFromFile()
	cacheEntries := GetAllFromCache()
	//zrobić taką samą pętlę jak w sync, i jak jest wpis w dump a nie ma w cache, to dodać do cache jako deleted
	diffs := make(map[string]*FileEntry)

	for _, dumpFile := range dumpEntries {
		found := false
		for _, cacheFile := range cacheEntries {
			if dumpFile.Key == cacheFile.Key {
				found = true
				break
			}
		}

		if !found {
			diffs[dumpFile.Key] = dumpFile
		}
	}

	for _, diff := range diffs {
		diff.Modified = time.Now().UTC()
		diff.Operation = Deleted
		CacheFile(diff.Key, diff)
	}
}
