package config

import "os"

type SinkConfig struct {
	Path         string `yaml:"path,omitempty" validate:"required"`
	Template     string `yaml:"template,omitempty"`
	TemplatePath string `yaml:"templatePath,omitempty"`
	Owner        string `yaml:"owner,omitempty" validate:"required_with=Group"`
	Group        string `yaml:"group,omitempty" validate:"required_with=Owner"`
	Mode         string `yaml:"mode,omitempty" validate:"fileMode"`

	// Hold update values when parsed
	UID          uint32
	GID          uint32
	FileMode     os.FileMode
}