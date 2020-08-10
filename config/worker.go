package config

import (
	"time"
)

type WorkerConfig struct {
	Resources     []ResourceConfig `yaml:"resources" validate:"required,dive,required"`
	Frequency     string           `yaml:"frequency,omitempty" validate:"required"`
	TimeFrequency time.Duration    `yaml:"timefrequency" validate:"-"`
	PreChange     string           `yaml:"preChange,omitempty"`
	PostChange    string           `yaml:"postChange,omitempty"`
	Sinks         []SinkConfig     `yaml:"sinks" validate:"required,dive,required"`
}
