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

type CyberarkResourceConfig struct {
	Alias        string       `yaml:"alias,omitempty"`
	Credential   string       `yaml:"credential,omitempty"`
	Name         string       `yaml:"name"`
	Version      string       `yaml:"version,omitempty"`
	Kind         ResourceKind `yaml:"kind,omitempty" validate:"required,oneof=cyberark-secret all-cyberark-secrets"`
	SafeName     string       `yaml:"safeName,omitempty" validate:"required"`
}

func (c CyberarkResourceConfig) GetName() string {
	return c.Name
}

func (c CyberarkResourceConfig) GetCredential() string {
	return c.Credential
}

func (c CyberarkResourceConfig) GetKind() ResourceKind {
	return c.Kind
}

func (c CyberarkResourceConfig) GetVault() string {
	return c.SafeName
}

func (c CyberarkResourceConfig) GetAlias() string {
	return c.Alias
}

func (c CyberarkResourceConfig) GetVersion() string {
	return c.Version
}
