package configparser

import (
	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

var validate *validator.Validate

type Config struct {
	Workers []config.WorkerConfig
}

func ParseConfig(path string) []config.WorkerConfig {
	config := Config{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Panicf("Error reading authconfig %v: %v", path, err)
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Panicf("Error unmarshalling yaml: %v", err)
	}

	validateWorkerConfigs(config.Workers)
	return config.Workers
}

func validateWorkerConfigs(workerConfigs []config.WorkerConfig) {
	validate = validator.New()

	for i, workerConfig := range workerConfigs {
		err := validate.Struct(workerConfig)

		// Convert human readable time and save into TimeFrequency
		workerConfigs[i].TimeFrequency = frequencyConverter(workerConfig.Frequency)
		if err != nil {
			log.Panicf("error: %v", err)
		}

		// Check each sinkConfig in the workerConfig
		for _, sinkConfig := range workerConfig.Sinks {
			// Ensure that Template and Template Path are not both defined
			if sinkConfig.Template != "" && sinkConfig.TemplatePath != "" {
				log.Panic("Template and TemplatePath cannot both be defined")
			}
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
