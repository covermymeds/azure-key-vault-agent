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
	for {
		select {
		case <-time.After(cfg.Frequency * time.Second):
			err := fetch(ctx, cfg)
			if err != nil {
				// TODO: If frequency > 1m, trigger a retry
				log.Printf("Failed to get resource: %v\n", err.Error())
			}

			//write(ctx, cfg)
		case <-ctx.Done():
			log.Println("Shutting down worker for: ", cfg.Name)
			return
		}
	}
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

	}

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
