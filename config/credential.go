package config

type CredentialConfig struct {
	Name         string `yaml:"name,omitempty" validate:"required"`
	TenantID     string `yaml:"tenantID,omitempty" validate:"required"`
	ClientID     string `yaml:"clientID,omitempty" validate:"required"`
	ClientSecret string `yaml:"clientSecret,omitempty" validate:"required"`
}
