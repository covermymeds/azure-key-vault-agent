package main

import (
	"flag"
	"github.com/chrisjohnson/azure-key-vault-agent/authconfig"
	"github.com/chrisjohnson/azure-key-vault-agent/configwatcher"
	log "github.com/sirupsen/logrus"
)

func init() {
	// JSON Format customized to use _timestamp so it marshals first alphabetically
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "_timestamp",
		},
	})

	var err error
	err = authconfig.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err.Error())
	}

	err = authconfig.AddFlags()
	if err != nil {
		log.Fatalf("failed to parse flags: %v\n", err.Error())
	}
	flag.Parse()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Caught Panic In Main: %v", r)
		}
	}()

	configwatcher.Watcher("akva.yaml")
}

