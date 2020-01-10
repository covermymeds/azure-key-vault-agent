package configwatcher

import (
	"context"
	"github.com/chrisjohnson/azure-key-vault-agent/configparser"
	"github.com/chrisjohnson/azure-key-vault-agent/sinkworker"
	"github.com/fsnotify/fsnotify"
	"log"
)

func Watcher(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error establishing file watcher: %v\n", err)
		return
	}

	// If something goes wrong along the way, close the watcher
	defer watcher.Close()

	done := make(chan bool)

	// Parse config and start workers.  Get the cancel function back so it can be passed to the file configwatcher
	cancel := parseAndStartWorkers(path)
	defer cancel()

	// Now that the workers have been started, watch the config file and bounce them if changes happen
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
	sinkConfigs := configparser.ParseConfig(path)
	for _, sinkConfig := range sinkConfigs {
		go sinkworker.Worker(ctx, sinkConfig)
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
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("Config watcher noticed a change to %v\n", event.Name)
				// Kill workers and start new ones
				cancel()
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
