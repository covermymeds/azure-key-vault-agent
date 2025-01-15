package config

type KeyvaultCredentialConfig struct {
	Name         string `yaml:"name,omitempty" validate:"required"`
	TenantID     string `yaml:"tenantID,omitempty" validate:"required"`
	ClientID     string `yaml:"clientID,omitempty" validate:"required"`
	ClientSecret string `yaml:"clientSecret,omitempty" validate:"required"`
}

func (kvc KeyvaultCredentialConfig) GetName() string {
	return kvc.Name
}
