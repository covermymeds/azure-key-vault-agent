package client

import (
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
)

type Client interface {
	GetCert(vault string, certName string, certVersion string) (certs.Cert, error)
	GetCerts(vault string) (results []certs.Cert, err error)
	GetSecret(vault string, secretName string, secretVersion string) (secrets.Secret, error)
	GetSecrets(vault string) (results map[string]secrets.Secret, err error)
	GetKey(vault string, keyName string, keyVersion string) (keys.Key, error)
	GetKeys(vault string) (results []keys.Key, err error)
}

type Clients map[string]Client
