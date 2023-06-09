package client

import (
	"fmt"

	"github.com/covermymeds/azure-key-vault-agent/config"
	"github.com/covermymeds/azure-key-vault-agent/iam"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
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

func NewSpnClient(cred config.CredentialConfig) keyvault.BaseClient {
	client := keyvault.New()
	authorizer, err := iam.GetKeyvaultAuthorizerFromSpn(cred.TenantID, cred.ClientID, cred.ClientSecret)
	if err != nil {
		panic(fmt.Sprintf("Error authorizing: %v", err.Error()))
	}
	client.Authorizer = autorest.AutorestAuthorizer(authorizer)
	return client
}

func NewCliClient() keyvault.BaseClient {
	client := keyvault.New()
	authorizer, err := iam.GetKeyvaultAuthorizerFromCli()
	if err != nil {
		panic(fmt.Sprintf("Error authorizing: %v", err.Error()))
	}
	client.Authorizer = autorest.AutorestAuthorizer(authorizer)
	return client
}
