package config

type ResourceKind string

const (
	CertKind       ResourceKind = "cert"
	KeyKind        ResourceKind = "key"
	SecretKind     ResourceKind = "secret"
	AllSecretsKind ResourceKind = "all-secrets"
)

type ResourceConfig struct {
	Kind         ResourceKind `yaml:"kind,omitempty" validate:"required,oneof=cert key secret"`
	VaultBaseURL string       `yaml:"vaultBaseURL,omitempty" validate:"required,url"`
	Name         string       `yaml:"name,omitempty" validate:"required"`
	Credential   string       `yaml:"credential,omitempty"`
	Version      string       `yaml:"version,omitempty"`
}
