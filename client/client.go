package client

import (
	"fmt"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/iam"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
)

type Clients map[string]keyvault.BaseClient

func NewClient(cred config.CredentialConfig) keyvault.BaseClient {
	client := keyvault.New()
	authorizer, err := iam.GetKeyvaultAuthorizer(cred.TenantID, cred.ClientID, cred.ClientSecret)
	if err != nil {
		panic(fmt.Sprintf("Error authorizing: %v", err.Error()))
	}
	client.Authorizer = authorizer
	return client
}
