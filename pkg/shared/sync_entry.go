package shared

import "time"

type SyncEntry struct {
	Path      string
	IsDir     bool
	Action    SyncAction
	Timestamp time.Time
}

type SyncAction int8

const (
	Upload = iota
	Download
	Remove
	Create
)
