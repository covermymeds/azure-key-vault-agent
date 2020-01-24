package configwatcher

import (
	"context"
	"github.com/chrisjohnson/azure-key-vault-agent/configparser"
	"github.com/chrisjohnson/azure-key-vault-agent/worker"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"time"
)

func Watcher(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicf("Error establishing file watcher: %v\n", err)
	}

	// If something goes wrong along the way, close the watcher
	defer watcher.Close()

	done := make(chan bool)

	// Parse authconfig and start workers.  Get the cancel function back so it can be passed to the file configwatcher
	cancel := parseAndStartWorkers(path)
	defer cancel()

	// Now that the workers have been started, watch the authconfig file and bounce them if changes happen
	go doWatch(watcher, cancel, path)

	err = watcher.Add(path)
	if err != nil {
		log.Panicf("Error watching path %v: %v\n", path, err)
	}
	<-done // Block until done
}

func parseAndStartWorkers(path string) context.CancelFunc {
	// Create background context for workers
	ctx, cancel := context.WithCancel(context.Background())

	// Parse config file and start workers
	workerConfigs := configparser.ParseConfig(path)
	for _, workerConfig := range workerConfigs {
		go worker.Worker(ctx, workerConfig)
	}

	return cancel
}

func doWatch(watcher *fsnotify.Watcher, cancel context.CancelFunc, path string) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
				log.Printf("Config watcher noticed a change to %v\n", event.Name)
				// Kill workers
				cancel()
				// Wait for file at path to be available (editors like Vim need to swap)
				waitForFile(path, 10)
				// Make a new Watcher
				Watcher(path)
				// Start new workers
				cancel = parseAndStartWorkers(path)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			log.Printf("Config watcher encountered an error for %v: %v\n", path, err)
			return
		}
	}
}

func waitForFile(path string, retries int){
	_, err := os.Stat(path)
	if retries == 0 {
		log.Panicf("Unable to find file: %v", path)
	}
	if os.IsNotExist(err) {
		log.Printf("Waiting for %v to be ready", path)
		time.Sleep(1*time.Second)
		waitForFile(path, retries-1)
	}
}
