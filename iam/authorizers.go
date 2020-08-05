// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package iam

import (
	"net/url"
	"strings"

	"github.com/chrisjohnson/azure-key-vault-agent/authconfig"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

// GetKeyvaultAuthorizer gets an OAuthTokenAuthorizer for use with Key Vault
// keys and secrets. Note that Key Vault *Vaults* are managed by Azure Resource
// Manager.
func GetKeyvaultAuthorizer() (autorest.Authorizer, error) {
	// BUG: default value for KeyVaultEndpoint is wrong
	vaultEndpoint := strings.TrimSuffix(authconfig.Environment().KeyVaultEndpoint, "/")
	// BUG: alternateEndpoint replaces other endpoints in the configs below
	alternateEndpoint, _ := url.Parse(
		"https://login.windows.net/" + authconfig.TenantID() + "/oauth2/token")

	var a autorest.Authorizer
	var err error

	oauthconfig, err := adal.NewOAuthConfig(
		authconfig.Environment().ActiveDirectoryEndpoint, authconfig.TenantID())
	if err != nil {
		return a, err
	}
	oauthconfig.AuthorizeEndpoint = *alternateEndpoint

	token, err := adal.NewServicePrincipalToken(
		*oauthconfig, authconfig.ClientID(), authconfig.ClientSecret(), vaultEndpoint)
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
