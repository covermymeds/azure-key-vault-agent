package keys

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/iam"
)

func getClient() keyvault.BaseClient {
	client := keyvault.New()
	a, err := iam.GetKeyvaultAuthorizer()
	if err != nil {
		log.Fatalf("Error authorizing: %v\n", err.Error())
	}
	client.Authorizer = a
	client.AddToUserAgent(config.UserAgent())
	return client
}

func GetKey(vaultBaseURL string, keyName string, keyVersion string) (result keyvault.JSONWebKey, err error) {
	client := getClient()

	key, err := client.GetKey(context.Background(), vaultBaseURL, keyName, keyVersion)
	if err != nil {
		log.Fatalf("Error getting key: %v\n", err.Error())
	}

	result = *key.Key

	return
}
