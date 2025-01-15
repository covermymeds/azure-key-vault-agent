package config

type CyberarkCredentialConfig struct {
	Name         string `yaml:"name,omitempty" validate:"required"`
	WorkloadID   string `yaml:"workloadID,omitempty" validate:"required"`
	ApiKey       string `yaml:"apiKey,omitempty" validate:"required"`
}

func (kvc CyberarkCredentialConfig) GetName() string {
	return kvc.Name
}
