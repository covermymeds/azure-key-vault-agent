package config

type CyberarkCredentialConfig struct {
	Name         string `yaml:"name,omitempty" validate:"required"`
	Login        string `yaml:"login,omitempty" validate:"required"`
	ApiKey       string `yaml:"apiKey,omitempty" validate:"required"`
	Account      string `yaml:"account,omitempty" validate:"required"`
	ApplianceURL string `yaml:"applianceURL,omitempty" validate:"required"`
}

func (kvc CyberarkCredentialConfig) GetName() string {
	return kvc.Name
}
