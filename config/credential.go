package config

type CredentialConfig struct {
	Name         string `yaml:"name,omitempty" validate:"required"`
	TenantID     string `yaml:"tenantID,omitempty" validate:"required_unless=CliAuth true"`
	ClientID     string `yaml:"clientID,omitempty" validate:"required_unless=CliAuth true"`
	ClientSecret string `yaml:"clientSecret,omitempty" validate:"required_unless=CliAuth true"`
	CliAuth      bool   `yaml:"cliAuth,omitempty"`
}
