package client

import (
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
)

type Client interface {
	GetCert(vaultBaseURL string, certName string, certVersion string) (certs.Cert, error)
	GetCerts(vaultBaseURL string) (results []certs.Cert, err error)
	GetSecret(vaultBaseURL string, secretName string, secretVersion string) (secrets.Secret, error)
	GetSecrets(vaultBaseURL string) (results map[string]secrets.Secret, err error)
	GetKey(vaultBaseURL string, keyName string, keyVersion string) (keys.Key, error)
	GetKeys(vaultBaseURL string) (results []keys.Key, err error)
}

type Clients map[string]Client
