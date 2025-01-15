package configwatcher

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/covermymeds/azure-key-vault-agent/client"
	"github.com/covermymeds/azure-key-vault-agent/config"
	"github.com/covermymeds/azure-key-vault-agent/configparser"
	"github.com/covermymeds/azure-key-vault-agent/worker"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func Watcher(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("Error establishing file watcher: %v", err))
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
		panic(fmt.Sprintf("Error watching path %v: %v", path, err))
	}
	<-done // Block until done
}

func initializeClients(parsedConfig configparser.Config) client.Clients{
	clients := make(client.Clients)
	for _, credentialConfig := range parsedConfig.Credentials {
		switch t := credentialConfig.CredConfig.(type) {
		case config.KeyvaultCredentialConfig:
			kvc := credentialConfig.CredConfig.(config.KeyvaultCredentialConfig)
			clients[credentialConfig.GetName()] = client.NewKeyvaultClient(kvc)
		case config.CyberarkCredentialConfig:
			panic("Unimplemented")
		default:
			panic(fmt.Sprintf("Got unexpected type: %v", t))
		}
	}

	return clients
}
func ParseAndRunWorkersOnce(path string) {
	// Parse config file
	parsedConfig := configparser.ParseConfig(path)

	// Initialize clients
	clients := initializeClients(parsedConfig)

	// Start workers
	log.Printf("Running workers once")
	for _, workerConfig := range parsedConfig.Workers {
		err := worker.Process(nil, clients, workerConfig)
		if err != nil {
			log.Fatalf("Failed to get resource(s): %v", err)
		}
	}
}

func parseAndStartWorkers(path string) context.CancelFunc {
	// Create background context for workers
	ctx, cancel := context.WithCancel(context.Background())

	// Parse config file
	parsedConfig := configparser.ParseConfig(path)

	// Initialize clients
	clients := initializeClients(parsedConfig)

	// Start workers
	for _, workerConfig := range parsedConfig.Workers {
		go worker.Worker(ctx, clients, workerConfig)
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
				log.Printf("Config watcher noticed a change to %v", event.Name)
				// Kill workers
				cancel()
				// Start new workers
				cancel = parseAndStartWorkers(path)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			log.Printf("Config watcher encountered an error for %v: %v", path, err)
			return
		}
	}
}
