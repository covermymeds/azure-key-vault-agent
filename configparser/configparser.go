package configparser

import (
	"github.com/chrisjohnson/azure-key-vault-agent/sink"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Sinks []sink.SinkConfig `yaml: "sinks, omitempty"`
}

func ParseConfig(path string) []sink.SinkConfig{
	ac := Config{}
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