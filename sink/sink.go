package sink

import (
	"time"
)

type SinkKind string

const (
	CertKind   SinkKind = "cert"
	KeyKind    SinkKind = "key"
	SecretKind SinkKind = "secret"
)

type SinkConfig struct {
	Kind         SinkKind      `yaml:"kind,omitempty" validate:"required"`
	VaultBaseURL string        `yaml:"vaultBaseURL,omitempty" validate:"required"`
	Name         string        `yaml:"name,omitempty" validate:"required"`
	Version      string        `yaml:"version,omitempty"`
	Path         string        `yaml:"path,omitempty" validate:"required"`
	Frequency    time.Duration `yaml:"frequency,omitempty"`
	Template     string        `yaml:"template,omitempty"`
	TemplatePath string        `yaml:"templatePath,omitempty"`
	PreChange    string        `yaml:"preChange,omitempty"`
	PostChange   string        `yaml:"postChange,omitempty"`

	/*
	   - kind: cert
	     path: /etc/nginx/certs/foo.cert
	     refresh: 5s
	     vaultBaseURL: https://cjohnson-kv.vault.azure.net/
	     name: cjohnson-test-cert
	     postChange: service nginx restart
	     preChange: who knows
	     version: latest # or specific version number
	*/
}
