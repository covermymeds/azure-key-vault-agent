package client

import (
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/config"

)

type CyberarkClient struct {}

func NewCyberarkClient(cred config.CyberarkCredentialConfig) CyberarkClient {
	return CyberarkClient{}
}

func (c CyberarkClient) GetCert(vaultBaseURL string, certName string, certVersion string) (certs.Cert, error) {
	panic("cyberark doesn't have Cert type resources. use regualar Secrets instead")
}

func (c CyberarkClient) GetCerts(vaultBaseURL string) (results []certs.Cert, err error) {
	panic("cyberark doesn't have Cert type resources. use regualar Secrets instead")
}

func (c CyberarkClient) GetSecret(vaultBaseURL string, secretName string, secretVersion string) (secrets.Secret, error) {
	v := "secretvaluefor" + secretName + "v" + secretVersion
	ct := "text"
	result := secrets.Secret{
		Value: &v,
		ContentType: &ct,
	}

	return result, nil
}

func (c CyberarkClient) GetSecrets(vaultBaseURL string) (results map[string]secrets.Secret, err error) {
	panic("unimplemented")
}

func (c CyberarkClient) GetKey(vaultBaseURL string, keyName string, keyVersion string) (keys.Key, error) {
	panic("cyberark does not have a Key secret type. use regular Secrets instead")
}
