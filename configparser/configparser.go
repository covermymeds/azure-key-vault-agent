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
	"regexp"
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

func ValidateFileMode(fl validator.FieldLevel) bool {
	// This is an optional field
	if fl.Field().String() == "" {
		return true
	}

	matched, err := regexp.MatchString(`^[0-7]{3,4}$`, fl.Field().String())
	if err != nil {
		panic(err)
	}

	return matched
}

func parseWorkerConfigs(workerConfigs []config.WorkerConfig) {
	validate = validator.New()
	validate.RegisterValidation("fileMode", ValidateFileMode)

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

	// Parse the Ownership
	sinkConfig = parseSinkOwnership(sinkConfig)

	// Parse the Permissions
	sinkConfig = parseSinkPermissions(sinkConfig)

	return sinkConfig
}

func parseSinkPermissions(sinkConfig config.SinkConfig) config.SinkConfig{
	if sinkConfig.Mode != "" {
		// Parse the last 3 digits for unix permissions
		permbits, err := strconv.ParseUint(sinkConfig.Mode[len(sinkConfig.Mode)-3:], 8, 32)
		if err != nil {
			panic(err)
		}

		// Set final mode to just perm bits for now
		finalMode := os.FileMode(permbits)

		// Calculate special bits if we have them i.e. 1700
		if len(sinkConfig.Mode) == 4 {
			// Get the Special bits
			specialBits, err := strconv.ParseUint(string(sinkConfig.Mode[0]), 8, 32)
			if err != nil {
				panic(err)
			}

			// Figure out if sticky, setgid, or setuid and apply proper bitwise or
			sticky := uint32(1)
			setgid := uint32(2)
			setuid := uint32(4)

			if uint32(specialBits)&sticky == sticky {
				finalMode = finalMode | os.ModeSticky
			}

			if uint32(specialBits)&setgid == setgid {
				finalMode = finalMode | os.ModeSetgid
			}

			if uint32(specialBits)&setuid == setuid {
				finalMode = finalMode | os.ModeSetuid
			}
		}

		// Update sinkConfig to reflect calculated file perms
		sinkConfig.FileMode = finalMode
	} else {
		// Set default file mode of 644
		sinkConfig.FileMode = os.FileMode(0644)
	}

	return sinkConfig
}

func parseSinkOwnership(sinkConfig config.SinkConfig) config.SinkConfig {
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
