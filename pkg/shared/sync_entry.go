package shared

import (
	proto_v1 "github.com/msoft-dev/filetify/pkg/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type SyncEntry struct {
	Key       string
	Path      string
	IsDir     bool // IsDir Consider if this is necessary, as this information should come in file transfer
	Action    SyncAction
	Timestamp time.Time //Timestamp Same like for IsDir
}

type SyncAction int8

const (
	Upload SyncAction = iota
	Download
	Remove
	CreateServer
	CreateLocal
)

func (e *SyncEntry) ToProto() *proto_v1.ServerFileSyncInfo {
	return &proto_v1.ServerFileSyncInfo{
		Key:       e.Key,
		Path:      e.Path,
		Action:    proto_v1.FileAction(e.Action + 1),
		Timestamp: timestamppb.New(e.Timestamp),
		IsDir:     e.IsDir,
	}
}

func ToSyncEntry(p *proto_v1.ServerFileSyncInfo) *SyncEntry {
	return &SyncEntry{
		Key:       p.Key,
		Path:      p.Path,
		IsDir:     p.IsDir,
		Timestamp: p.Timestamp.AsTime(),
		Action:    SyncAction(p.Action - 1),
	}
}
