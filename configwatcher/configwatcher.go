package configwatcher

import (
	"context"
	"fmt"
	"github.com/chrisjohnson/azure-key-vault-agent/configparser"
	"github.com/chrisjohnson/azure-key-vault-agent/worker"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func Watcher(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("Error establishing file watcher: %v\n", err))
	}

	// If something goes wrong along the way, close the watcher
	defer watcher.Close()

	done := make(chan bool)

	// Parse authconfig and start workers.  Get the cancel function back so it can be passed to the file configwatcher
	cancel := parseAndStartWorkers(path)
	defer cancel()

	// Now that the workers have been started, watch the authconfig file and bounce them if changes happen
	go doWatch(watcher, cancel, path)

	err = watcher.Add(filepath.Dir(path))
	if err != nil {
		panic(fmt.Sprintf("Error watching path %v: %v\n", path, err))
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
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Caught Panic In doWatch: %v", r)
		}
	}()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			if (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) && event.Name == path {
				log.Printf("Config watcher noticed a change to %v\n", event.Name)
				// Kill workers
				cancel()
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
