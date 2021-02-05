package main

import (
	"flag"
	"fmt"
	"github.com/covermymeds/azure-key-vault-agent/configwatcher"
	"github.com/luci/luci-go/common/flag/flagenum"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

type outputType uint

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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
var debugMode bool
var ver bool
var runOnce bool

func init() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.SetOutput(os.Stdout)

	fs.BoolVar(&help, "help", false, "Show this help text")
	fs.BoolVar(&debugMode, "debug", false, "Enable debugging")
	fs.StringVar(&configFile, "config", "", "Read config from this `file`")
	fs.StringVar(&configFile, "c", "", "Read config from this `file` (shorthand)")
	fs.Var(&output, "output", fmt.Sprintf("Output type (default json). Options are: %v (default json)", outputTypeEnum.Choices()))
	fs.BoolVar(&ver, "version", false, "Show the version of akva-key-vault-agent")
	fs.BoolVar(&runOnce, "once", false, "Run once and quit")

	fs.Parse(os.Args[1:])

	if help {
		fs.PrintDefaults()
		os.Exit(0)
	}

	if ver {
		fmt.Printf("Version: %s\nCommit: %s\nBuilt on %s\n", version, commit, date)
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

	fs.Parse(os.Args[1:])
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if debugMode {
				debug.PrintStack()
			}
			log.Fatalf("Caught panic in main: %v", r)
		}
	}()

	if runOnce {
		configwatcher.ParseAndRunWorkersOnce(configFile)
	} else {
		configwatcher.Watcher(configFile)
	}
}
