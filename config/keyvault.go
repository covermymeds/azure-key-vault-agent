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

type KeyvaultResourceConfig struct {
	Alias        string       `yaml:"alias,omitempty"`
	Credential   string       `yaml:"credential,omitempty"`
	Name         string       `yaml:"name"`
	Version      string       `yaml:"version,omitempty"`
	Kind         ResourceKind `yaml:"kind,omitempty" validate:"required,oneof=cert key secret all-secrets"`
	VaultBaseURL string       `yaml:"vaultBaseURL,omitempty" validate:"required,url"`
}

func (k KeyvaultResourceConfig) GetName() string {
	return k.Name
}

func (k KeyvaultResourceConfig) GetCredential() string {
	return k.Credential
}

func (k KeyvaultResourceConfig) GetKind() ResourceKind {
	return k.Kind
}

func (k KeyvaultResourceConfig) GetVault() string {
	return k.VaultBaseURL
}

func (k KeyvaultResourceConfig) GetAlias() string {
	return k.Alias
}

func (k KeyvaultResourceConfig) GetVersion() string {
	return k.Version
}
