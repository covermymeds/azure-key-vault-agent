package config

type ResourceKind string

const (
	CertKind       ResourceKind = "cert"
	KeyKind        ResourceKind = "key"
	SecretKind     ResourceKind = "secret"
	AllSecretsKind ResourceKind = "all-secrets"
)

type ResourceConfig struct {
	Alias        string       `yaml:"alias,omitempty"`
	Credential   string       `yaml:"credential,omitempty"`
	Kind         ResourceKind `yaml:"kind,omitempty" validate:"required,oneof=cert key secret all-secrets"`
	Name         string       `yaml:"name"`
	VaultBaseURL string       `yaml:"vaultBaseURL,omitempty" validate:"required,url"`
	Version      string       `yaml:"version,omitempty"`
}
