package client

import (
	"errors"
	"github.com/msoft-dev/filetify/pkg/shared"
)

func Synchronize() {
	conn, ok := GetConn()
	if !ok {
		shared.HandleError(errors.New("cannot get client connection"), false)
		return
	}

	defer conn.Close()

	filesToSync := StartSynchronization(conn)

	var uploads []*shared.SyncEntry
	var downloads []*shared.SyncEntry

	//TODO: removes, clientAdd

	for _, entry := range filesToSync {
		switch entry.Action {
		case shared.Upload:
			uploads = append(uploads, entry)
		case shared.Download:
			downloads = append(downloads, entry)
		}

	}

	//Upload files
	for _, file := range uploads {
		UploadFile(conn, file)
	}

	//Download files

	//TODO: After synchronization remove deleted files from cache
}
