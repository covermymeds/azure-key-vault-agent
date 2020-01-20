package config

type SinkConfig struct {
	Path         string `yaml:"path,omitempty" validate:"required"`
	Template     string `yaml:"template,omitempty"`
	TemplatePath string `yaml:"templatePath,omitempty"`
}