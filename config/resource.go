package config

type ResourceKind string

const (
	CertKind       ResourceKind = "cert"
	KeyKind        ResourceKind = "key"
	SecretKind     ResourceKind = "secret"
	AllSecretsKind ResourceKind = "all-secrets"
)

type ResourceConfig struct {
	Kind         ResourceKind `yaml:"kind,omitempty" validate:"required,oneof=cert key secret all-secrets"`
	VaultBaseURL string       `yaml:"vaultBaseURL,omitempty" validate:"required,url"`
	Name         string       `yaml:"name"`
	Credential   string       `yaml:"credential,omitempty"`
	Version      string       `yaml:"version,omitempty"`
}
