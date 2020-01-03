package main

import (
	"flag"
	"log"

	"github.com/chrisjohnson/azure-key-vault-agent/certs"
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
}
