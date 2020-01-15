package configparser

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/chrisjohnson/azure-key-vault-agent/sink"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
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

	for i, sinkConfig := range sinkConfigs {
		err := validate.Struct(sinkConfig)
		sinkConfigs[i].TimeFrequency = frequencyConverter(sinkConfig.Frequency)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

func frequencyConverter(freq string) time.Duration {
	readabletime, _ := time.ParseDuration(freq)

	// If the time specified is less than 1 second, treat the value as seconds
	if readabletime <= time.Duration(time.Second) {
		intreadable, err := strconv.Atoi(freq)
		if err != nil {
			// Default time to 1m instead of 100ms if input is not valid
			readabletime, _ = time.ParseDuration("1m")
			log.Println("The value of Frequency was not valid, defaulting to 1m from", freq)
		} else {
			// Convert the nanoseconds to seconds
			readabletime = time.Duration(intreadable) * time.Second
			log.Println("The value of Frequency was too low, defaulting to seconds from", freq)

		}

	}
	return readabletime
}
