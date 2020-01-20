package sink

import "time"

type SinkKind string

const (
	CertKind   SinkKind = "cert"
	KeyKind    SinkKind = "key"
	SecretKind SinkKind = "secret"
)

type ConfigFile struct {
	Workers []WorkerConfig `yaml:"workers,omitempty" validate:"required"`
}

type WorkerConfig struct {
	Resources     []ResourceConfig `yaml:"resources,omitempty" validate:"required"`
	Frequency     string           `yaml:"frequency,omitempty" validate:"required"`
	TimeFrequency time.Duration    `yaml:"timefrequency" validate:"-"`
	PreChange     string           `yaml:"preChange,omitempty"`
	PostChange    string           `yaml:"postChange,omitempty"`
	Sinks         []SinkConfig     `yaml:"sinks,omitempty" validate:"required"`
}

type ResourceConfig struct {
	Kind         SinkKind `yaml:"kind,omitempty" validate:"required,oneof=cert key secret"`
	VaultBaseURL string   `yaml:"vaultBaseURL,omitempty" validate:"required,url"`
	Name         string   `yaml:"name,omitempty" validate:"required"`
	Version      string   `yaml:"version,omitempty"`
}

type SinkConfig struct {
	Path         string `yaml:"path,omitempty" validate:"required"`
	Template     string `yaml:"template,omitempty"`
	TemplatePath string `yaml:"templatePath,omitempty"`
}
