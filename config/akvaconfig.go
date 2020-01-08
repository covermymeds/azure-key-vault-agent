package config

import (
	"github.com/chrisjohnson/azure-key-vault-agent/sink"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type AkvaConfig struct {
	Sinks []sink.SinkConfig `yaml: "sinks, omitempty"`
}

func ParseAkvaConfig(path string) []sink.SinkConfig{
	ac := AkvaConfig{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &ac)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return ac.Sinks
}