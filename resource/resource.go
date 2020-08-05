package resource

import (
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
)

type Resource interface {
	String() string
}

type ResourceMap struct {
	Certs   map[string]certs.Cert
	Secrets map[string]secrets.Secret
	Keys    map[string]keys.Key
}
