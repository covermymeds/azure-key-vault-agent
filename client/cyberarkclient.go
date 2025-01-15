package client

import (
	"fmt"

	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/config"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/secrets"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

type CyberarkClient struct {
	Client *conjurapi.Client
}

func NewCyberarkClient(cred config.CyberarkCredentialConfig) CyberarkClient {
	cyberarkConfig := conjurapi.Config{
		Account: cred.Account,
		ApplianceURL: cred.ApplianceURL,
	}

	cyberarkClient, err := conjurapi.NewClientFromKey(cyberarkConfig,
		authn.LoginPair{
			Login: cred.Login,
			APIKey: cred.ApiKey,
		},
	)
	if err != nil {
		panic(err)
	}
	return CyberarkClient{cyberarkClient}
}

func (c CyberarkClient) GetCert(vaultBaseURL string, certName string, certVersion string) (certs.Cert, error) {
	panic("cyberark doesn't have Cert type resources. use regular Secrets instead")
}

func (c CyberarkClient) GetCerts(vaultBaseURL string) (results []certs.Cert, err error) {
	panic("cyberark doesn't have Cert type resources. use regular Secrets instead")
}

func (c CyberarkClient) GetSecret(vaultBaseURL string, secretName string, secretVersion string) (secrets.Secret, error) {
	secretValue, err := c.Client.RetrieveSecret(secretName)
	if err != nil {
			panic(err)
	}

	secretValueString := string(secretValue)
	result := secrets.Secret{
		Value: &secretValueString,
		ContentType: nil,
	}

	return result, nil
}

func (c CyberarkClient) GetSecrets(vaultBaseURL string) (results map[string]secrets.Secret, err error) {
	resources, err := c.Client.ResourceIDs(&conjurapi.ResourceFilter{Kind: "variable"})
	if err != nil {
		panic(err)
	}
	secretValues, err := c.Client.RetrieveBatchSecrets(resources)
	if err != nil {
		panic(err)
	}

	results = make(map[string]secrets.Secret)

	for resourceID, value := range secretValues {
		secretValueString := string(value)
		result := secrets.Secret{
			Value: &secretValueString,
			ContentType: nil,
		}

		fmt.Printf("%s = '%s'\n", resourceID, secretValueString)
		results[resourceID] = result
	}

	return results, nil
}

func (c CyberarkClient) GetKey(vaultBaseURL string, keyName string, keyVersion string) (keys.Key, error) {
	panic("cyberark does not have a Key secret type. use regular Secrets instead")
}
