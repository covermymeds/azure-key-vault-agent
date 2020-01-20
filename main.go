package main

import (
	"flag"
	"github.com/chrisjohnson/azure-key-vault-agent/authconfig"
	"github.com/chrisjohnson/azure-key-vault-agent/configwatcher"
	"log"
)

func init() {
	var err error
	err = authconfig.ParseEnvironment()
	if err != nil {
		log.Panicf("failed to parse env: %v\n", err.Error())
	}

	err = authconfig.AddFlags()
	if err != nil {
		log.Panicf("failed to parse flags: %v\n", err.Error())
	}
	flag.Parse()
}

func main() {
	configwatcher.Watcher("akva.yaml")
}
