package main

import (
	"flag"
	"log"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
)

func init() {
	var err error
	err = config.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err.Error())
	}

	err = config.AddFlags()
	if err != nil {
		log.Fatalf("failed to parse flags: %v\n", err.Error())
	}
	flag.Parse()
}

func main() {
	// vault url, secret name, version (can leave blank for "latest")
	secret, err := secrets.GetSecret("https://cjohnson-kv.vault.azure.net/", "password", "8f1e2267024a4dacb81b14aa33b8f084")
	if err != nil {
		log.Fatalf("failed to get secret: %v\n", err.Error())
	}
	log.Println(secret)

	secrets, listErr := secrets.GetSecrets("https://cjohnson-kv.vault.azure.net/")
	if listErr != nil {
		log.Fatalf("failed to get list of secrets: %v\n", listErr.Error())
	}
	for _, value := range secrets {
		log.Println(value)
	}
}
