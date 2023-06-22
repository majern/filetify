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
	updateCacheFromDump()

	//Start FileWatcher
	startFileWatcher()

	//Schedule Synchronization
	go scheduleSynchronization()

	//Schedule cache dump
	go scheduleDumpCacheToFile()
	select {}
}

func scanFiles() {
	shared.ScanFiles(GetConfiguration().Paths, false, GetConfiguration().IgnoredFiles)
}

func updateCacheFromDump() {
	shared.UpdateCacheFromDump()
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

// scheduleDumpCacheToFile dumps the cache to the file
// Deprecated: TODO: it should be called when close app signal comes from the system
func scheduleDumpCacheToFile() {
	var interval = GetConfiguration().SyncTimeoutSec
	for {
		time.Sleep(time.Duration(interval*10) * time.Second)
		go shared.DumpToFile()
	}
}
