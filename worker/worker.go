package worker

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/certs"
	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/keys"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
	"github.com/chrisjohnson/azure-key-vault-agent/templaterenderer"

	"github.com/jpillora/backoff"
)

const RetryBreakPoint = 60

func Worker(ctx context.Context, workerConfig config.WorkerConfig) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Caught Panic In Worker: %v", r)
		}
	}()
	b := &backoff.Backoff{
		Min:    time.Duration(workerConfig.TimeFrequency),
		Max:    time.Duration(workerConfig.TimeFrequency) * 10,
		Factor: 3,
		Jitter: true,
	}

	d := b.Duration()
	ticker := time.NewTicker(d)

	log.Printf("Starting worker with frequency %v", d)

	err := Process(ctx, workerConfig)
	if err != nil {
		log.Printf("Failed to get resource(s): %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			// The main thread has cancelled the worker
			log.Println("Shutting down worker")
			return
		case <-ticker.C:
			err := Process(ctx, workerConfig)
			if err != nil {
				if workerConfig.TimeFrequency > RetryBreakPoint {
					// For long frequencies, we should set up an explicit retry (with backoff)
					// For short frequencies, we can just wait for the next natural tick
					d = b.Duration()
					ticker = time.NewTicker(d)
				}
				log.Println(err)
				log.Printf("Failed to get resource(s), will retry in %v", d)
			} else {
				// Reset the ticker once we've got a good result
				if workerConfig.TimeFrequency > RetryBreakPoint {
					b.Reset()
					d = b.Duration()
					ticker = time.NewTicker(d)
				}
				log.Printf("Successfully fetched resource(s), will try next in %v", d)
			}
		}
	}
}

func Process(ctx context.Context, workerConfig config.WorkerConfig) error {
	resources := resource.ResourceMap{make(map[string]certs.Cert), make(map[string]secrets.Secret), make(map[string]keys.Key)}

	for _, resourceConfig := range workerConfig.Resources {
		switch resourceConfig.Kind {
		case config.CertKind:
			result, err := certs.GetCert(resourceConfig.VaultBaseURL, resourceConfig.Name, resourceConfig.Version)
			if err != nil {
				return err
			}
			resources.Certs[resourceConfig.Name] = result

		case config.SecretKind:
			result, err := secrets.GetSecret(resourceConfig.VaultBaseURL, resourceConfig.Name, resourceConfig.Version)
			if err != nil {
				return err
			}
			resources.Secrets[resourceConfig.Name] = result

		case config.KeyKind:
			result, err := keys.GetKey(resourceConfig.VaultBaseURL, resourceConfig.Name, resourceConfig.Version)
			if err != nil {
				return err
			}
			resources.Keys[resourceConfig.Name] = result

		default:
			panic(fmt.Sprintf("Invalid sink kind: %v", resourceConfig.Kind))
		}
	}

	type Change struct {
		sinkConfig  config.SinkConfig
		newContents string
	}

	var changes []Change
	for _, sinkConfig := range workerConfig.Sinks {
		// Get old content
		oldContents := getOldContent(sinkConfig)

		// Get new content
		newContents := getNewContent(sinkConfig, resources)

		// If a change was detected run pre/post commands and write the new file
		if oldContents != newContents {
			changes = append(changes, Change{sinkConfig, newContents})
			log.Printf("Change detected for %v", sinkConfig.Path)
		}
	}

	if len(changes) > 0 {
		if workerConfig.PreChange != "" {
			err := runCommand(workerConfig.PreChange)
			if err != nil {
				log.Printf("PreChange command errored: %v", err)
			}
		}

		for _, change := range changes {
			write(change.sinkConfig, change.newContents)
		}

		if workerConfig.PostChange != "" {
			err := runCommand(workerConfig.PostChange)
			if err != nil {
				log.Printf("PostChange command errored: %v", err)
			}
		}
	}

	return nil
}

func fetch(ctx context.Context, resourceConfig config.ResourceConfig) (result resource.Resource, err error) {
	switch resourceConfig.Kind {
	case config.CertKind:
		result, err = certs.GetCert(resourceConfig.VaultBaseURL, resourceConfig.Name, resourceConfig.Version)

	case config.SecretKind:
		result, err = secrets.GetSecret(resourceConfig.VaultBaseURL, resourceConfig.Name, resourceConfig.Version)

	case config.KeyKind:
		result, err = keys.GetKey(resourceConfig.VaultBaseURL, resourceConfig.Name, resourceConfig.Version)

	default:
		panic(fmt.Sprintf("Invalid sink kind: %v", resourceConfig.Kind))
	}

	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func getNewContent(sinkConfig config.SinkConfig, resources resource.ResourceMap) string {
	// If we have templates get the new value from rendering them
	if sinkConfig.Template != "" || sinkConfig.TemplatePath != "" {
		if sinkConfig.Template != "" {
			// Execute inline template
			return templaterenderer.RenderInline(sinkConfig.Template, resources)
		} else {
			// Execute template file
			return templaterenderer.RenderFile(sinkConfig.TemplatePath, resources)
		}
	} else {
		// Just return the string
		// TODO: If there is only one resource being requested, call .String() on it
		return "TODO"
	}
}

func getOldContent(sinkConfig config.SinkConfig) string {
	// If path has changed it will not yet exist so return empty string
	if _, err := os.Stat(sinkConfig.Path); err != nil {
		if os.IsNotExist(err) {
			return ""
		}
	}

	// Read the contents of the current file into a string
	b, err := ioutil.ReadFile(sinkConfig.Path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func write(sinkConfig config.SinkConfig, content string) {
	f, err := os.Create(sinkConfig.Path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(content)

	if err != nil {
		panic(err)
	}
}

func runCommand(command string) error {
	log.Printf("Executing %v", command)
	cmd := exec.Command("sh", "-c", command)

	stdoutStderr, err := cmd.CombinedOutput()
	if stdoutStderr != nil {
		log.Printf(string(stdoutStderr))
	}

	return err
}
