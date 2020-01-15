package configparser

import (
	"github.com/chrisjohnson/azure-key-vault-agent/sink"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var validate *validator.Validate

type Config struct {
	Sinks []sink.SinkConfig `yaml: "sinks, omitempty"`
}

func ParseConfig(path string) []sink.SinkConfig {
	config := Config{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Panicf("Error reading config %v: %v", path, err)
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Panicf("Error unmarshalling yaml: %v", err)
	}

	validateSinkConfigs(config.Sinks)

	return config.Sinks
}

func validateSinkConfigs(sinkConfigs []sink.SinkConfig) {
	validate = validator.New()

	for _, sinkConfig := range sinkConfigs {
		err := validate.Struct(sinkConfig)

		if err != nil {
			log.Panicf("error: %v", err)
		}

		// Ensure that Template and Template Path are not both defined
		if sinkConfig.Template != "" && sinkConfig.TemplatePath != "" {
			log.Panic("Template and TemplatePath cannot both be defined")
		}
	}
}
