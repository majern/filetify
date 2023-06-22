package shared

import (
	"github.com/fsnotify/fsnotify"
	proto_v1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type FileEntry struct {
	Key       string
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

func NewFileEntry(key, path string, isDir bool, modTime time.Time, op FileOperation) *FileEntry {
	entry := &FileEntry{Key: key, Path: path, IsDir: isDir, Modified: modTime, Operation: op}
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

func (e *FileEntry) ToProto() *proto_v1.FileSyncInfo {
	return &proto_v1.FileSyncInfo{
		Key:       e.Key,
		Path:      e.Path,
		Status:    proto_v1.FileStatus(e.Operation + 1),
		Timestamp: timestamppb.New(e.Modified),
		IsDir:     e.IsDir,
	}
}

func ToFileEntry(p *proto_v1.FileSyncInfo) *FileEntry {
	return &FileEntry{
		Key:       p.Key,
		Path:      p.Path,
		IsDir:     p.IsDir,
		Modified:  p.Timestamp.AsTime(),
		Operation: FileOperation(p.Status - 1),
	}
}
