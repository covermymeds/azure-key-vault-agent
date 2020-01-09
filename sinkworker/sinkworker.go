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
	log.Println("Starting worker for: ", cfg.Name, "with refresh: ", cfg.Frequency)

	interval := cfg.Frequency
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	retry := false

	process(ctx, cfg)
	err := process(ctx, cfg)
	if err != nil {
		log.Printf("Failed to get resource: %v\n", err.Error())
	}

	for {
		log.Printf("Polling for worker %v\n", cfg.Name)
		select {
		case <-ctx.Done():
			// The main thread has cancelled the worker
			log.Println("Shutting down worker for: ", cfg.Name)
			return
		case <-ticker.C:
			err := process(ctx, cfg)
			if err != nil {
				if cfg.Frequency > 60 {
					// For long frequencies, we should set up an explicit retry
					// Shorter frequencies, we can just wait for the next natural tick
					if retry {
						// Double the next retry interval
						interval = interval * 2
						ticker = time.NewTicker(time.Duration(interval) * time.Second)
					} else {
						// First failure, reset interval to 1 and enable retry logic
						interval = 1
						retry = true
					}
				}
				log.Printf("Failed to get resource: %v\n", err.Error())
			} else {
				// Reset the ticker once we've got a good result
				if cfg.Frequency > 60 {
					retry = false
					interval = cfg.Frequency
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
	log.Println("Fetching:", cfg.Name)
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

/*
	// vault url, secret name, version (can leave blank for "latest")
	secret, err := secrets.GetSecret("https://cjohnson-kv.vault.azure.net/", "password", "8f1e2267024a4dacb81b14aa33b8f084")
	if err != nil {
		log.Fatalf("failed to get secret: %v\n", err.Error())
	}
	log.Printf("Got secret password: %v\n", secret)

	secrets, listErr := secrets.GetSecrets("https://cjohnson-kv.vault.azure.net/")
	if listErr != nil {
		log.Fatalf("failed to get list of secrets: %v\n", listErr.Error())
	}
	log.Println("Getting all secrets")
	for _, value := range secrets {
		log.Println(value)
	}
	log.Println("Done")

	// vault url, cert name, version (can leave blank for "latest")
	cert, err := certs.GetCert("https://cjohnson-kv.vault.azure.net/", "cjohnson-test", "4cffd52057214a0799287e2ea905ffd9")
	if err != nil {
		log.Fatalf("failed to get cert: %v\n", err.Error())
	}
	log.Printf("Got cert cjohnson-test: %v\n", cert)

	certs, listErr := certs.GetCerts("https://cjohnson-kv.vault.azure.net/")
	if listErr != nil {
		log.Fatalf("failed to get list of certs: %v\n", listErr.Error())
	}
	log.Println("Getting all certs")
	for _, value := range certs {
		log.Println(value)
	}
	log.Println("Done")
*/
