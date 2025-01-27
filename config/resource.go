package config

import "fmt"

type ResourceKind string

const (
	CertKind               ResourceKind = "cert"
	KeyKind                ResourceKind = "key"
	SecretKind             ResourceKind = "secret"
	AllSecretsKind         ResourceKind = "all-secrets"
	AllCyberarkSecretsKind ResourceKind = "all-cyberark-secrets"
	CyberarkSecretKind     ResourceKind = "cyberark-secret"
)

type GenericResource interface {
	GetName()       string
	GetCredential() string
	GetKind()       ResourceKind
	GetVault()      string
	GetVersion()    string
	GetAlias()      string
}

type ResourceConfig struct {
	GenericResource
}

func (c *ResourceConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var kvResource KeyvaultResourceConfig
	if err := unmarshal(&kvResource); err == nil && kvResource.VaultBaseURL != "" {
		if kvResource.Credential == "" {
			kvResource.Credential = "default"
		}
		c.GenericResource = kvResource
		return nil
	}

	var caResource CyberarkResourceConfig
	if err := unmarshal(&caResource); err == nil && caResource.SafeName != "" {
		if caResource.Credential == "" {
			caResource.Credential = "default_cyberark"
		}
		c.GenericResource = caResource
		return nil
	}

	return fmt.Errorf("unrecognized resource type")
}
