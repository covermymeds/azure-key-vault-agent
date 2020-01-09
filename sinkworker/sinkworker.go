package sinkworker

import (
	"context"
	"log"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/certs"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
	"github.com/chrisjohnson/azure-key-vault-agent/sink"
)

func Worker(ctx context.Context, cfg sink.SinkConfig) {
	retry := false
	interval := time.Duration(cfg.Frequency * time.Second)
	ticker := time.NewTicker(interval)

	log.Println("Starting worker for: ", cfg.Name, "with refresh: ", interval)

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
				if cfg.Frequency > 60 {
					// For long frequencies, we should set up an explicit retry
					// Shorter frequencies, we can just wait for the next natural tick
					if retry {
						// Double the next retry interval
						interval = time.Duration(interval * 2)
						ticker = time.NewTicker(interval)
						log.Printf("Failed to get resource %v, backing off and will retry in %v\n%v\n", cfg.Name, interval, err.Error())
					} else {
						// First failure, reset interval to 30 and enable retry logic
						interval = 30
						retry = true
						log.Printf("Failed to get resource %v, will retry in %v\n%v\n", cfg.Name, interval, err.Error())
					}
				} else {
					log.Printf("Failed to get resource: %v\n", err.Error())
				}
			} else {
				// Reset the ticker once we've got a good result
				if cfg.Frequency > 60 {
					retry = false
					interval = cfg.Frequency * time.Second
					ticker = time.NewTicker(time.Duration(interval) * time.Second)
				}
			}
		}
	}
}

func process(ctx context.Context, cfg sink.SinkConfig) (err error) {
	err = fetch(ctx, cfg)
	if err == nil {
		return
	}

	err = write(ctx, cfg)
	if err == nil {
		return
	}

	return
}

func fetch(ctx context.Context, cfg sink.SinkConfig) (err error) {
	switch cfg.Kind {
	case sink.CertKind:
		cert, certErr := certs.GetCert(cfg.VaultBaseURL, cfg.Name, cfg.Version)
		if certErr != nil {
			err = certErr
		}
		log.Printf("Got cert %v: %v\n", cfg.Name, cert)
		// TODO: Send to file writer, along with any template details
		// TODO: Trigger pre and post change hooks
		// TODO: Determine what constitutes a "change"
		// TODO: retry/backoff
		// TODO: If freq < 1m, just ignore failures, otherwise create an explicit retry

	case sink.SecretKind:
		secret, secretErr := secrets.GetSecret(cfg.VaultBaseURL, cfg.Name, cfg.Version)
		if secretErr != nil {
			err = secretErr
		}
		log.Printf("Got secret %v: %v\n", cfg.Name, secret)

	case sink.KeyKind:
		log.Fatalf("Not implemented yet")

	default:
		log.Fatalf("Invalid sink kind: %v\n", cfg.Kind)
	}

	return
}

func write(ctx context.Context, cfg sink.SinkConfig) (err error) {
	err = nil
	return
}
