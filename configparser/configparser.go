package configparser

import (
	"fmt"
	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
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
		panic(fmt.Sprintf("Error reading config %v: %v", path, err))
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling yaml: %v", err))
	}

	parseWorkerConfigs(config.Workers)
	return config.Workers
}

func parseWorkerConfigs(workerConfigs []config.WorkerConfig) {
	validate = validator.New()

	for i, workerConfig := range workerConfigs {
		err := validate.Struct(workerConfig)
		if err != nil {
			panic(fmt.Sprintf("Error parsing worker config: %v", err))
		}

		// Convert human readable time and save into TimeFrequency
		workerConfigs[i].TimeFrequency = frequencyConverter(workerConfig.Frequency)

		// Check each sinkConfig in the workerConfig
		for j, sinkConfig := range workerConfig.Sinks {
			workerConfigs[i].Sinks[j] = parseSinkConfig(sinkConfig)
		}
	}
}

func parseSinkConfig(sinkConfig config.SinkConfig) config.SinkConfig {
	// Ensure that Template and Template Path are not both defined
	if sinkConfig.Template != "" && sinkConfig.TemplatePath != "" {
		panic("Template and TemplatePath cannot both be defined")
	}

	//Check the Owner and Group for existence
	if sinkConfig.Owner != "" {
		u, err := user.Lookup(sinkConfig.Owner)
		if err != nil {
			panic(err)
		}

		uid, err := strconv.ParseUint(u.Uid, 10, 32)
		if err != nil {
			panic(err)
		}

		sinkConfig.UID = uint32(uid)
	} else {
		// Default to calling UID
		sinkConfig.UID = uint32(os.Getuid())
	}

	if sinkConfig.Group != "" {
		g, err := user.LookupGroup(sinkConfig.Group)
		if err != nil {
			panic(err)
		}

		gid, err := strconv.ParseUint(g.Gid, 10, 32)
		if err != nil {
			panic(err)
		}

		sinkConfig.GID = uint32(gid)
	} else {
		// Default to calling GID
		sinkConfig.GID = uint32(os.Getgid())
	}

	return sinkConfig
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
