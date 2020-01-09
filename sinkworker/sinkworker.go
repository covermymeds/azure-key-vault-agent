package sinkworker

import (
	"context"
	"log"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/certs"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
	"github.com/chrisjohnson/azure-key-vault-agent/sink"

	"github.com/jpillora/backoff"
)

const RetryBreakPoint = 60

func Worker(ctx context.Context, cfg sink.SinkConfig) {
	b := &backoff.Backoff{
		Min:    cfg.Frequency * time.Second,
		Max:    cfg.Frequency * 10 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	d := b.Duration()
	ticker := time.NewTicker(d)

	log.Printf("Starting worker for %v with refresh %v\n", cfg.Name, d)

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
				if cfg.Frequency > RetryBreakPoint {
					// For long frequencies, we should set up an explicit retry
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
				if cfg.Frequency > RetryBreakPoint {
					b.Reset()
					d := b.Duration()
					ticker = time.NewTicker(d)
					log.Printf("Success for resource %v, will try next in %v\n", cfg.Name, d)
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

	/*
		err = write(ctx, cfg)
		if err == nil {
			return
		}
	*/

	return
}

func fetch(ctx context.Context, cfg sink.SinkConfig) (err error) {
	switch cfg.Kind {
	case sink.CertKind:
		cert, certErr := certs.GetCert(cfg.VaultBaseURL, cfg.Name, cfg.Version)
		if certErr != nil {
			err = certErr
		} else {
			log.Printf("Got cert %v: %v\n", cfg.Name, cert)
		}
		// TODO: Send to file writer, along with any template details
		// TODO: Trigger pre and post change hooks
		// TODO: Determine what constitutes a "change"
		// TODO: retry/backoff
		// TODO: If freq < 1m, just ignore failures, otherwise create an explicit retry

	case sink.SecretKind:
		secret, secretErr := secrets.GetSecret(cfg.VaultBaseURL, cfg.Name, cfg.Version)
		if secretErr != nil {
			err = secretErr
		} else {
			log.Printf("Got secret %v: %v\n", cfg.Name, secret)
		}

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
