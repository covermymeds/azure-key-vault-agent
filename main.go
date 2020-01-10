package main

import (
	"flag"
	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/configwatcher"
	"log"
)

func init() {
	var err error
	err = config.ParseEnvironment()
	if err != nil {
		log.Panicf("failed to parse env: %v\n", err.Error())
	}

	err = config.AddFlags()
	if err != nil {
		log.Panicf("failed to parse flags: %v\n", err.Error())
	}
	flag.Parse()
}

func main() {
	configwatcher.Watcher("akva.yaml")
}
