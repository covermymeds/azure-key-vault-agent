package keys

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/iam"
)

func getKeysClient() keyvault.BaseClient {
	keyClient := keyvault.New()
	a, err := iam.GetKeyvaultAuthorizer()
	if err != nil {
		log.Fatalf(err.Error())
	}
	keyClient.Authorizer = a
	keyClient.AddToUserAgent(config.UserAgent())
	return keyClient
}

func GetSecret(vaultBaseURL string, secretName string, secretVersion string) (result string, err error) {
	keysClient := getKeysClient()

	secret, err := keysClient.GetSecret(context.Background(), vaultBaseURL, secretName, secretVersion)
	if err != nil {
		log.Fatalf("Error getting secret: %v\n", err.Error())
	}

	result = *secret.Value

	return
}
