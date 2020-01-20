package config

import (
	"time"
)

type WorkerConfig struct {
	Resources     []ResourceConfig
	Frequency     string            `yaml:"frequency,omitempty" validate:"required"`
	TimeFrequency time.Duration     `yaml:"timefrequency" validate:"-"`
	PreChange     string            `yaml:"preChange,omitempty"`
	PostChange    string            `yaml:"postChange,omitempty"`
	Sinks         []SinkConfig      `yaml: "sinks, omitempty"`
}

