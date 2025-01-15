package config

import (
	"fmt"
)

type CredConfig interface {
	GetName() string
}

type CredentialConfig struct {
	CredConfig
	CredentialType string
}

type KeyvaultCredentialConfig struct {
	Name           string `yaml:"name,omitempty" validate:"required"`
	TenantID     string `yaml:"tenantID,omitempty" validate:"required"`
	ClientID     string `yaml:"clientID,omitempty" validate:"required"`
	ClientSecret string `yaml:"clientSecret,omitempty" validate:"required"`
}

func (kvc KeyvaultCredentialConfig) GetName() string {
	return kvc.Name
}

type CyberarkCredentialConfig struct {
	Name           string `yaml:"name,omitempty" validate:"required"`
	WorkloadID   string `yaml:"workloadID,omitempty" validate:"required"`
	ApiKey       string `yaml:"apiKey,omitempty" validate:"required"`
}

func (kvc CyberarkCredentialConfig) GetName() string {
	return kvc.Name
}

func (c *CredentialConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var kvConfig KeyvaultCredentialConfig
	if err := unmarshal(&kvConfig); err == nil && kvConfig.ClientID != "" {
		c.CredConfig = kvConfig
		c.CredentialType = "azure-keyvault"
		return nil
	}

	var caConfig CyberarkCredentialConfig
	if err := unmarshal(&caConfig); err == nil && caConfig.WorkloadID != "" {
		c.CredConfig = caConfig
		c.CredentialType = "cyberark-conjur"
		return nil
	}

	return fmt.Errorf("unrecognized credential type")
}
