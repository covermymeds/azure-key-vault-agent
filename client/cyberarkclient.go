package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/config"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
	log "github.com/sirupsen/logrus"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

type CyberarkClient struct {
	Client *conjurapi.Client
	Safe string
}

func NewCyberarkClient(cred config.CyberarkCredentialConfig) CyberarkClient {
	cyberarkConfig := conjurapi.Config{
		Account: cred.Account,
		ApplianceURL: cred.ApplianceURL,
	}

	cyberarkClient, err := conjurapi.NewClientFromKey(cyberarkConfig,
		authn.LoginPair{
			Login: fmt.Sprintf("host/data/%s", cred.Login),
			APIKey: cred.ApiKey,
		},
	)
	if err != nil {
		panic(fmt.Sprintf("Error creating Cyberark client: %v", err.Error()))
	}
	return CyberarkClient{Client: cyberarkClient, Safe: cred.Safe}
}

func (c CyberarkClient) GetCert(vaultBaseURL string, certName string, certVersion string) (certs.Cert, error) {
	panic("cyberark doesn't have Cert type resources. use regular Secrets instead")
}

func (c CyberarkClient) GetCerts(vaultBaseURL string) (results []certs.Cert, err error) {
	panic("cyberark doesn't have Cert type resources. use regular Secrets instead")
}

func (c CyberarkClient) GetSecret(vaultBaseURL string, secretName string, secretVersion string) (secrets.Secret, error) {
	var secretValue []byte
	var err error

	secretPath := fmt.Sprintf("data/vault/%s/%s", c.Safe, secretName)

	if secretVersion == "" {
		secretValue, err = c.Client.RetrieveSecret(secretPath)
	} else {
		secretVersionInt, convErr := strconv.Atoi(secretVersion)
		if convErr != nil {
			return secrets.Secret{}, fmt.Errorf("failed to convert secret version to integer: %s", secretVersion)
		}
		secretValue, err = c.Client.RetrieveSecretWithVersion(secretPath, secretVersionInt)
	}
	if err != nil {
		log.Printf("Error getting secret: %v", err.Error())
		return secrets.Secret{}, err
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
		log.Printf("Error getting secrets: %v", err.Error())
		return map[string]secrets.Secret{}, err
	}

	secretValues, err := c.Client.RetrieveBatchSecrets(resources)
	if err != nil {
		log.Printf("Error getting secrets: %v", err.Error())
		return map[string]secrets.Secret{}, err
	}

	results = make(map[string]secrets.Secret)

	for resourceID, value := range secretValues {
		modResourceID := strings.Replace(resourceID, fmt.Sprintf("conjur:variable:data/vault/%s/", c.Safe), "", 1)
		secretValueString := string(value)
		result := secrets.Secret{
			Value: &secretValueString,
			ContentType: nil,
		}

		results[modResourceID] = result
	}

	return results, nil
}

func (c CyberarkClient) GetKey(vaultBaseURL string, keyName string, keyVersion string) (keys.Key, error) {
	panic("cyberark does not have a Key secret type. use regular Secrets instead")
}

func (c CyberarkClient) GetKeys(vaultBaseURL string) ([]keys.Key, error) {
	panic("cyberark does not have a Key secret type. use regular Secrets instead")
}
