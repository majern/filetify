package server

import (
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
)

func SyncFiles(remoteFiles []*shared.FileEntry) []*shared.SyncEntry {
	logrus.WithField("files", remoteFiles).Debug("Syncing files...")
	localFiles := shared.GetAllFromCache()
	syncEntries := compareFiles(remoteFiles, localFiles)

	logrus.WithField("filesToSync", syncEntries).Debug("Found files to sync")

	return syncEntries
}

// compareFiles: Compare files and returns slice of synchronization actions that should be performed
// on client side.
func compareFiles(remote, local []*shared.FileEntry) []*shared.SyncEntry {
	diffs := make(map[string]*shared.SyncEntry)

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, remoteFile := range remote {
			found := false
			for _, localFile := range local {
				if remoteFile.Path == localFile.Path {
					found = true
					if diffs[remoteFile.Path] == nil {
						compRes := compareFile(remoteFile, localFile)
						switch compRes {
						case -1:
							diffs[remoteFile.Path] = &shared.SyncEntry{Key: remoteFile.Key, Path: remoteFile.Path, Action: shared.Upload, IsDir: remoteFile.IsDir, Timestamp: remoteFile.Modified}
						case 1:
							diffs[remoteFile.Path] = &shared.SyncEntry{Key: remoteFile.Key, Path: remoteFile.Path, Action: shared.Download, IsDir: remoteFile.IsDir, Timestamp: remoteFile.Modified}
						}
					}
					break
				}
			}
			// Path not found. We add it to return slice
			if !found && diffs[remoteFile.Path] == nil {
				if i == 0 {
					diffs[remoteFile.Path] = &shared.SyncEntry{Key: remoteFile.Key, Path: remoteFile.Path, Action: shared.Upload, IsDir: remoteFile.IsDir, Timestamp: remoteFile.Modified}
				} else {
					diffs[remoteFile.Path] = &shared.SyncEntry{Key: remoteFile.Key, Path: remoteFile.Path, Action: shared.Download, IsDir: remoteFile.IsDir, Timestamp: remoteFile.Modified}
				}
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			remote, local = local, remote
		}
	}

	result := make([]*shared.SyncEntry, 0, len(diffs))

	for _, v := range diffs {
		result = append(result, v)
	}

	return result
}

// compareFile: Compare remote and local file. If remote is newer than local then
// it returns -1, if both are the same it returns 0. If local is newer, it returns 1.
func compareFile(remote, local *shared.FileEntry) int8 {
	if remote.Modified.Before(local.Modified) {
		return 1
	}
	if remote.Modified.After(local.Modified) {
		return -1
	}
	return 0
}
