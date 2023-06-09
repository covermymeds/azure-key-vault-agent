package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func GetKeyvaultAuthorizerFromSpn(tenantID string, clientID string, clientSecret string) (auth.Authorizer, error) {
	environment := environments.AzurePublic()
	credentials := auth.Credentials{
		Environment:                           *environment,
		EnableAuthenticatingUsingClientSecret: true,
		TenantID:                              tenantID,
		ClientID:                              clientID,
		ClientSecret:                          clientSecret,
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.KeyVault)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}

	return authorizer, err
}
