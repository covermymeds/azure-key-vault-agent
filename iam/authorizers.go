// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package iam

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	environment *azure.Environment
)

const CloudName string = "AzurePublicCloud"

// GetKeyvaultAuthorizer gets an OAuthTokenAuthorizer for use with Key Vault
// keys and secrets. Note that Key Vault *Vaults* are managed by Azure Resource
// Manager.
func GetKeyvaultAuthorizer(tenantID string, clientID string, clientSecret string) (autorest.Authorizer, error) {
	// BUG: default value for KeyVaultEndpoint is wrong
	vaultEndpoint := strings.TrimSuffix(getEnvironment().KeyVaultEndpoint, "/")
	// BUG: alternateEndpoint replaces other endpoints in the configs below
	alternateEndpoint, _ := url.Parse(
		"https://login.windows.net/" + tenantID + "/oauth2/token")

	var a autorest.Authorizer

	oauthconfig, err := adal.NewOAuthConfig(
		getEnvironment().ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return a, err
	}
	oauthconfig.AuthorizeEndpoint = *alternateEndpoint

	token, err := adal.NewServicePrincipalToken(
		*oauthconfig, clientID, clientSecret, vaultEndpoint)
	if err != nil {
		return a, err
	}

	a = autorest.NewBearerAuthorizer(token)

	if err == nil {
		return a, nil
	} else {
		return nil, err
	}
}

func getEnvironment() *azure.Environment {
	if environment != nil {
		return environment
	}
	env, err := azure.EnvironmentFromName(CloudName)
	if err != nil {
		panic(fmt.Sprintf(
			"invalid cloud name '%s' specified, cannot continue", CloudName))
	}
	environment = &env
	return environment
}
