package shared

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

const dir = "/Users/mariusz/OwnCloud/Projects/MSoft/Repositories"

var watcher *fsnotify.Watcher

func read() {
	var computerName, _ = os.Hostname()
	println(computerName)

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	fmt.Println("Adding directories to file watcher...")

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(dir, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	fmt.Println("Direcotires added!")

	//
	done := make(chan bool)
	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if !strings.Contains(event.Name, "~") {
					switch event.Op {
					case fsnotify.Write:
						fmt.Println("File UPDATED", event.Name, " --> SEND NEW")
					case fsnotify.Create:
						fmt.Println("File CREATED", event.Name, " --> SEND NEW")
					case fsnotify.Remove:
						fmt.Println("File REMOVED", event.Name, " --> DELETE")
					case fsnotify.Rename:
						fmt.Println("File RENAMED", event.Name, " --> DELETE")
					}
				}
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		var err = watcher.Add(path)

		if err != nil {
			fmt.Printf("An error occurred while adding %v: %v\n", path, err)
		} else {
			fmt.Printf("Path added to watcher: %v\n", path)
		}

		return err
	}

	return nil
}
