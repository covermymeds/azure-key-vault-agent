package resource

import (
	"github.com/chrisjohnson/azure-key-vault-agent/certs"
	"github.com/chrisjohnson/azure-key-vault-agent/keys"
	"github.com/chrisjohnson/azure-key-vault-agent/secrets"
)

type Resource interface {
	String() string
}

type ResourceMap struct {
	Certs   map[string]certs.Cert
	Secrets map[string]secrets.Secret
	Keys    map[string]keys.Key
}
