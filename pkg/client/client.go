package client

import (
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"time"
)

func StartClient() {
	shared.InitLogger(&GetConfiguration().LogConfig)
	//TODO: This should be rewrited to use https://github.com/kardianos/service and https://github.com/madflojo/tasks
	logrus.Info("Starting client")

	scanFiles()

	//Start FileWatcher
	startFileWatcher()

	//Schedule Synchronization
	go scheduleSynchronization()
	select {}
}

func scanFiles() {
	shared.ScanFiles(GetConfiguration().Paths, GetConfiguration().IgnoredFiles)
}

func startFileWatcher() {
	paths := GetConfiguration().Paths

	for _, path := range paths {
		go WatchFilesInPath(path)
	}
}

func scheduleSynchronization() {
	var interval = GetConfiguration().SyncTimeoutSec
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		go Synchronize()
	}
}
