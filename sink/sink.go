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
	Kind         SinkKind      `json:"kind,omitempty"`
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
