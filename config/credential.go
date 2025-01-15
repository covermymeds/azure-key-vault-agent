package config

import (
	"fmt"
)

type CredConfig interface {
	GetName() string
}

type CredentialConfig struct {
	CredConfig
}

func (c *CredentialConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var kvConfig KeyvaultCredentialConfig
	if err := unmarshal(&kvConfig); err == nil && kvConfig.ClientID != "" {
		c.CredConfig = kvConfig
		return nil
	}

	var caConfig CyberarkCredentialConfig
	if err := unmarshal(&caConfig); err == nil && caConfig.WorkloadID != "" {
		c.CredConfig = caConfig
		return nil
	}

	return fmt.Errorf("unrecognized credential type")
}
