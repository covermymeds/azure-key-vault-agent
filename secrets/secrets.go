package secrets

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"

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

func GetSecret(vaultBaseURL string, secretName string, secretVersion string) (result string, err error) {
	client := getClient()

	secret, err := client.GetSecret(context.Background(), vaultBaseURL, secretName, secretVersion)
	if err != nil {
		log.Fatalf("Error getting secret: %v\n", err.Error())
	}

	result = *secret.Value

	return
}

func GetSecretByURL(secretURL string) (result string, err error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		log.Fatalf("Failed to parse URL for secret: %v\n", err.Error())
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/secrets/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	secretName := res[0][1]

	result, err = GetSecret(vaultBaseURL, secretName, "")
	if err != nil {
		log.Fatalf("Failed to get secret from parsed values %v and %v: %v\n", vaultBaseURL, secretName, err.Error())
	}

	return
}

func GetSecrets(vaultBaseURL string) (results []string, err error) {
	client := getClient()

	max := int32(25)
	slrp, err := client.GetSecrets(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Fatalf("Error getting secret: %v\n", err.Error())
	}

	for _, value := range slrp.Values() {
		secretURL := *value.ID
		secret, err := GetSecretByURL(secretURL)
		if err != nil {
			log.Fatalf("Error loading secret contents: %v\n", err.Error())
		}
		results = append(results, secret)
	}

	return
}
