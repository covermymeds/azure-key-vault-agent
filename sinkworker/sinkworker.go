package sinkworker

import (
	"context"
	"errors"
	"github.com/chrisjohnson/azure-key-vault-agent/templateparser"
	"log"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/certs"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
	"github.com/chrisjohnson/azure-key-vault-agent/sink"

	"github.com/jpillora/backoff"
)

const RetryBreakPoint = 60

func Worker(ctx context.Context, cfg sink.SinkConfig) {
	b := &backoff.Backoff{
		Min:    time.Duration(cfg.TimeFrequency),
		Max:    time.Duration(cfg.TimeFrequency) * 10,
		Factor: 2,
		Jitter: true,
	}

	d := b.Duration()
	ticker := time.NewTicker(d)

	log.Printf("Starting worker of kind %v for %v with refresh %v\n", cfg.Kind, cfg.Name, d)

	err := process(ctx, cfg)
	if err != nil {
		log.Printf("Failed to get resource: %v\n", err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			// The main thread has cancelled the worker
			log.Println("Shutting down worker for: ", cfg.Name)
			return
		case <-ticker.C:
			log.Printf("Polling for worker %v\n", cfg.Name)
			err := process(ctx, cfg)
			if err != nil {
				if cfg.TimeFrequency > RetryBreakPoint {
					// For long frequencies, we should set up an explicit retry (with backoff)
					d := b.Duration()
					ticker = time.NewTicker(d)
					log.Println(err)
					log.Printf("Failed to get resource %v, will retry in %v\n", cfg.Name, d)
				} else {
					// For short frequencies, we can just wait for the next natural tick
					log.Printf("Failed to get resource: %v\n", err.Error())
				}
			} else {
				// Reset the ticker once we've got a good result
				if cfg.TimeFrequency > RetryBreakPoint {
					b.Reset()
					d := b.Duration()
					ticker = time.NewTicker(d)
					log.Printf("Success for resource %v, will try next in %v\n", cfg.Name, d)
				}
			}
		}
	}
}

var count int

func process(ctx context.Context, cfg sink.SinkConfig) error {
	result, err := fetch(ctx, cfg)
	count++
	if count > 2 && count < 8 {
		return errors.New("FAKE ERROR FROM AZURE")
	}
	if err != nil {
		return err
	}
	log.Print(result.Map())
	log.Printf("Got resource of kind %v for %v: %v\n", cfg.Kind, cfg.Name, result.String())

	// TODO: return now if the value hasn't changed

	if cfg.Template != "" || cfg.TemplatePath != "" {
		log.Println("TODO: pass to template")
		if cfg.Template != "" {
			// execute inline template
			templateparser.InlineTemplate(cfg.Template, cfg.Path, result)
		} else {
			// execute template file
			templateparser.TemplateFile(cfg.TemplatePath, cfg.Path, result)
		}
	}

	if cfg.PreChange != "" {
		log.Println("TODO: prechange hooks")
	}

	err = write(ctx, cfg, result)
	if err != nil {
		return err
	}

	if cfg.PostChange != "" {
		log.Println("TODO: postchange hooks")
	}

	return nil
}

func fetch(ctx context.Context, cfg sink.SinkConfig) (result resource.Resource, err error) {
	switch cfg.Kind {
	case sink.CertKind:
		result, err = certs.GetCert(cfg.VaultBaseURL, cfg.Name, cfg.Version)

	case sink.SecretKind:
		result, err = secrets.GetSecret(cfg.VaultBaseURL, cfg.Name, cfg.Version)

	case sink.KeyKind:
		log.Panicf("Not implemented yet")

	default:
		log.Panicf("Invalid sink kind: %v\n", cfg.Kind)
	}

	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func write(ctx context.Context, cfg sink.SinkConfig, result resource.Resource) error {
	return nil
}
