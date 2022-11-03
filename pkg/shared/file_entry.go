package shared

import (
	"github.com/fsnotify/fsnotify"
	"time"
)

type FileEntry struct {
	Path      string
	IsDir     bool
	Modified  time.Time
	Operation FileOperation
	//Files     []*FileEntry
}

type FileOperation int8

const (
	Scanned FileOperation = iota
	New
	Modified
	Deleted
)

func NewFileEntry(path string, isDir bool, modTime time.Time, op FileOperation) *FileEntry {
	entry := &FileEntry{Path: path, IsDir: isDir, Modified: modTime, Operation: op}
	return entry
}

func ParseOperation(op fsnotify.Op) FileOperation {
	switch op {
	case fsnotify.Create:
		return New
	case fsnotify.Write:
		return Modified
	case fsnotify.Remove, fsnotify.Rename:
		return Deleted
	default:
		return Scanned
	}
}
