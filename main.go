package main

import (
	"log"
	"flag"

	"azure-key-vault-agent/config"
	"azure-key-vault-agent/keys"
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
	// secret name, version (can leave blank for "latest")
	secret, err := keys.GetSecret("cjohnson-kv", "password", "8f1e2267024a4dacb81b14aa33b8f084")
	if err != nil {
		log.Fatalf("failed to get secret: %v\n", err.Error())
	}

	log.Println(secret)
}
