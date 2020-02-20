package main

import (
	"flag"
	"fmt"
	"github.com/chrisjohnson/azure-key-vault-agent/authconfig"
	"github.com/chrisjohnson/azure-key-vault-agent/configwatcher"
	"github.com/luci/luci-go/common/flag/flagenum"
	log "github.com/sirupsen/logrus"
	"os"
)

type outputType uint

var outputTypeEnum = flagenum.Enum{
	"json": outputType(10),
	"text": outputType(20),
}

func (val *outputType) Set(v string) error {
	return outputTypeEnum.FlagSet(val, v)
}

func (val *outputType) String() string {
	return outputTypeEnum.FlagString(*val)
}

func (val outputType) MarshalJSON() ([]byte, error) {
	return outputTypeEnum.JSONMarshal(val)
}

var configFile string
var output outputType
var help bool

func init() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.SetOutput(os.Stdout)

	fs.StringVar(&configFile, "config", "", "Read config from this file")
	fs.StringVar(&configFile, "c", "", "Read config from this file (shorthand)")
	fs.BoolVar(&help, "help", false, "Show this help text")
	fs.Var(&output, "output", fmt.Sprintf("Output type (default json). Options are: %v (default json)", outputTypeEnum.Choices()))

	fs.Parse(os.Args[1:])

	if help {
		fs.PrintDefaults()
		os.Exit(0)
	}

	if configFile == "" {
		log.Fatalf("Missing --config/-c")
	}

	if output != outputTypeEnum["text"] {
		// JSON Format customized to use _timestamp so it marshals first alphabetically
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime: "_timestamp",
			},
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	var err error
	err = authconfig.ParseEnvironment()
	if err != nil {
		log.Fatalf("AuthConfig: Failed to parse env: %v", err.Error())
	}

	err = authconfig.AddFlags(*fs)
	if err != nil {
		log.Fatalf("AuthConfig: Failed to add flags: %v", err.Error())
	}

	fs.Parse(os.Args[1:])
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Caught Panic In Main: %v", r)
		}
	}()

	configwatcher.Watcher(configFile)
}
