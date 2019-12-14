package vaults

import (
	"fmt"
	"context"
	"log"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	//"github.com/Azure/go-autorest/autorest/to"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
)

func getVaultsClient() keyvault.VaultsClient {
	// create a Vaults client
	client := keyvault.NewVaultsClient(config.SubscriptionID())

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		client.Authorizer = authorizer
	}

	return client
}

// GetVault returns an existing vault
func GetVault(ctx context.Context, vaultName string) (keyvault.Vault, error) {
	vaultsClient := getVaultsClient()
	return vaultsClient.Get(ctx, config.ResourceGroupName(), vaultName)
}

func GetVaults() {
	vaultsClient := getVaultsClient()

	fmt.Println("Getting all vaults in subscription")
	for subList, err := vaultsClient.ListComplete(context.Background(), nil); subList.NotDone(); err = subList.Next() {
		if err != nil {
			log.Printf("failed to get list of vaults: %v", err)
		}
		fmt.Printf("\t%s\n", *subList.Value().Name)
	}

	fmt.Println("Getting all vaults in resource group")
	for rgList, err := vaultsClient.ListByResourceGroupComplete(context.Background(), config.ResourceGroupName(), nil); rgList.NotDone(); err = rgList.Next() {
		if err != nil {
			log.Printf("failed to get list of vaults: %v", err)
		}
		fmt.Printf("\t%s\n", *rgList.Value().Name)
	}
}
