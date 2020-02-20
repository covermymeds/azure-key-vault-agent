package main

import (
	"flag"
	"github.com/chrisjohnson/azure-key-vault-agent/authconfig"
	"github.com/chrisjohnson/azure-key-vault-agent/configwatcher"
	log "github.com/sirupsen/logrus"
)

var configFile string
var jsonLogs bool = true

func init() {
	var textLogs bool
	flag.BoolVar(&textLogs, "text-logs", false, "Change from JSON log output format to text")
	flag.StringVar(&configFile, "config", "", "Read config from this file")
	flag.StringVar(&configFile, "c", "", "Read config from this file")

	flag.Parse()

	if configFile == "" {
		log.Fatalf("Missing --config/-c")
	}

	if textLogs {
		jsonLogs = false
	}

	if jsonLogs {
		// JSON Format customized to use _timestamp so it marshals first alphabetically
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime: "_timestamp",
			},
		})
	}

	var err error
	err = authconfig.ParseEnvironment()
	if err != nil {
		log.Fatalf("AuthConfig: Failed to parse env: %v", err.Error())
	}

	err = authconfig.AddFlags()
	if err != nil {
		log.Fatalf("AuthConfig: Failed to add flags: %v", err.Error())
	}

	flag.Parse()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Caught Panic In Main: %v", r)
		}
	}()

	configwatcher.Watcher(configFile)
}
