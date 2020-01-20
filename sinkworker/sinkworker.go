package sinkworker

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/certs"
	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/keys"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
	"github.com/chrisjohnson/azure-key-vault-agent/templateparser"

	"github.com/jpillora/backoff"
)

const RetryBreakPoint = 60

func Worker(ctx context.Context, cfg config.WorkerConfig) {
	b := &backoff.Backoff{
		Min:    time.Duration(cfg.TimeFrequency),
		Max:    time.Duration(cfg.TimeFrequency) * 10,
		Factor: 2,
		Jitter: true,
	}

	d := b.Duration()
	ticker := time.NewTicker(d)

	log.Printf("Starting worker with refresh %v\n", d)

	err := process(ctx, cfg)
	if err != nil {
		log.Printf("Failed to get resource: %v\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			// The main thread has cancelled the worker
			log.Println("Shutting down worker")
			return
		case <-ticker.C:
			err := process(ctx, cfg)
			if err != nil {
				if cfg.TimeFrequency > RetryBreakPoint {
					// For long frequencies, we should set up an explicit retry (with backoff)
					d := b.Duration()
					ticker = time.NewTicker(d)
					log.Println(err)
					log.Printf("Failed to get resource(s), will retry in %v\n", d)
				} else {
					// For short frequencies, we can just wait for the next natural tick
					log.Printf("Failed to get resource: %v\n", err)
				}
			} else {
				// Reset the ticker once we've got a good result
				if cfg.TimeFrequency > RetryBreakPoint {
					b.Reset()
					d := b.Duration()
					ticker = time.NewTicker(d)
				}
				log.Printf("Success for resource(s), will try next in %v\n", d)
			}
		}
	}
}

func process(ctx context.Context, cfg config.WorkerConfig) error {
	resources := resource.ResourceMap{make(map[string]certs.Cert), make(map[string]secrets.Secret), make(map[string]keys.Key)}

	for _, resourceConfig := range cfg.Resources {
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
			log.Panicf("Invalid sink kind: %v\n", resourceConfig.Kind)
		}
	}

	for _, sinkConfig := range cfg.Sinks {
		// Get old content
		oldContent := getOldContent(sinkConfig)

		// Get new content
		newContent := getNewContent(sinkConfig, resources)

		// If a change was detected run pre/post commands and write the new file
		if oldContent != newContent {
			if cfg.PreChange != "" {
				err := runCommand(cfg.PreChange)
				if err != nil {
					log.Printf("PreChange command errored: %v", err)
				}
			}

			write(sinkConfig, newContent)

			if cfg.PostChange != "" {
				err := runCommand(cfg.PostChange)
				if err != nil {
					log.Printf("PostChange command errored: %v", err)
				}
			}
		}
	}

	return nil
}

func fetch(ctx context.Context, cfg config.ResourceConfig) (result resource.Resource, err error) {
	switch cfg.Kind {
	case config.CertKind:
		result, err = certs.GetCert(cfg.VaultBaseURL, cfg.Name, cfg.Version)

	case config.SecretKind:
		result, err = secrets.GetSecret(cfg.VaultBaseURL, cfg.Name, cfg.Version)

	case config.KeyKind:
		result, err = keys.GetKey(cfg.VaultBaseURL, cfg.Name, cfg.Version)

	default:
		log.Panicf("Invalid sink kind: %v\n", cfg.Kind)
	}

	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func getNewContent(cfg config.SinkConfig, resources resource.ResourceMap) string {
	// If we have templates get the new value from rendering them
	if cfg.Template != "" || cfg.TemplatePath != "" {
		if cfg.Template != "" {
			// Execute inline template
			return templateparser.InlineTemplate(cfg.Template, cfg.Path, resources)
		} else {
			// Execute template file
			return templateparser.TemplateFile(cfg.TemplatePath, cfg.Path, resources)
		}
	} else {
		// Just return the string
		return "TODO"
	}
}

func getOldContent(cfg config.SinkConfig) string {
	// If path has changed it will not yet exist so return empty string
	if _, err := os.Stat(cfg.Path); err != nil {
		if os.IsNotExist(err) {
			return ""
		}
	}

	// Read the contents of the current file into a string
	b, err := ioutil.ReadFile(cfg.Path)
	if err != nil {
		log.Panic(err)
	}

	return string(b)
}

func write(cfg config.SinkConfig, content string) {
	f, err := os.Create(cfg.Path)

	if err != nil {
		log.Panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(content)

	if err != nil {
		log.Panic(err)
	}
}

func runCommand(command string) error {
	log.Printf("Executing %v", command)
	cmd := exec.Command("sh", "-c", command)

	err := cmd.Run()
	return err
}
