package sink

import (
	"time"
)

type SinkType int

const (
	CertType SinkType = iota
	KeyType
	SecretType
)

type SinkConfig struct {
	Type         SinkType      `json:"type,omitempty"`
	VaultBaseURL string        `json:"vaultBaseURL,omitempty"`
	Name         string        `json:"name,omitempty"`
	Version      string        `json:"version,omitempty"`
	Path         string        `json:"path,omitempty"`
	Frequency    time.Duration `json:"value,omitempty"`
	Template     string        `json:"template,omitempty"`
	TemplatePath string        `json:"templatePath,omitempty"`
	PreChange    string        `json:"preChange,omitempty"`
	PostChange   string        `json:"postChange,omitempty"`

	/*
	   - type: cert
	     path: /etc/nginx/certs/foo.cert
	     refresh: 5s
	     vaultName: cjohnson-kv-test
	     name: cjohnson-test-cert
	     postChange: service nginx restart
	     preChange: who knows
	     version: latest # or specific version number
	*/
}
