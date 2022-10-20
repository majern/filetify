package client

import (
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/sirupsen/logrus"
	"time"
)

func StartClient() {
	shared.InitLogger()
	//TODO: This should be rewrited to use https://github.com/kardianos/service and https://github.com/madflojo/tasks
	logrus.WithTime(time.Now()).Info("Starting client")

	scanFiles()

	//Start FileWatcher
	startFileWatcher()

	//Schedule Synchronization
	go scheduleSynchronization()
	select {}
}

func scanFiles() {
	shared.ScanFiles(shared.GetConfiguration().Paths)
}

func startFileWatcher() {
	paths := shared.GetConfiguration().Paths

	for _, path := range paths {
		go WatchFilesInPath(path)
	}
}

func scheduleSynchronization() {
	var interval = shared.GetConfiguration().SyncTimeoutSec
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		logrus.Info("File synchronization started")
		shared.GetAllFromCache()
	}
}
