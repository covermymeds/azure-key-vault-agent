package keys

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"

	"azure-key-vault-agent/iam"
	"azure-key-vault-agent/config"
	"azure-key-vault-agent/vaults"
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

func GetSecret(vaultName string, secretName string, secretVersion string) (result string, err error) {
	keysClient := getKeysClient()

	vault, err := vaults.GetVault(context.Background(), vaultName)
	if err != nil {
		log.Fatalf("Failure getting vault: %v\n", err.Error())
	}

	vaultBaseURL := *vault.Properties.VaultURI

	secret, err := keysClient.GetSecret(context.Background(), vaultBaseURL, secretName, secretVersion)
	if err != nil {
		log.Fatalf("Error getting secret: %v\n", err.Error())
	}

	result = *secret.Value

	return
}
