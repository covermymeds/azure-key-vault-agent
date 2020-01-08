package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/sink"
	"github.com/chrisjohnson/azure-key-vault-agent/sinkworker"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// TODO: Parse config once to determine these
	// TODO: Set up a file watch on the config that parses, cancels the existing workers, and starts a new set

	cfg1 := sink.SinkConfig{Name: "username", Frequency: 1, Type: sink.SecretType, VaultBaseURL: "https://cjohnson-kv.vault.azure.net/"}
	cfg2 := sink.SinkConfig{Name: "password", Frequency: 1, Type: sink.SecretType, VaultBaseURL: "https://cjohnson-kv.vault.azure.net/"}
	cfg3 := sink.SinkConfig{Name: "cjohnson-test", Frequency: 1, Type: sink.CertType, VaultBaseURL: "https://cjohnson-kv.vault.azure.net/"}

	go sinkworker.Worker(ctx, cfg1)
	go sinkworker.Worker(ctx, cfg2)
	go sinkworker.Worker(ctx, cfg3)

	// Run for 15 seconds, then stop
	time.Sleep(15 * time.Second)
	log.Println("Shutting down")
	cancel()
}
