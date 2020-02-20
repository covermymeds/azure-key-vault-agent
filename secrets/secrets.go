package secrets

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"

	"github.com/chrisjohnson/azure-key-vault-agent/authconfig"
	"github.com/chrisjohnson/azure-key-vault-agent/iam"
)

type Secret keyvault.SecretBundle

func (s Secret) String() string {
	return *s.Value
}

func newClient() keyvault.BaseClient {
	client := keyvault.New()
	a, err := iam.GetKeyvaultAuthorizer()
	if err != nil {
		panic(fmt.Sprintf("Error authorizing: %v", err.Error()))
	}
	client.Authorizer = a
	client.AddToUserAgent(authconfig.UserAgent())
	return client
}

func GetSecret(vaultBaseURL string, secretName string, secretVersion string) (Secret, error) {
	secret, err := newClient().GetSecret(context.Background(), vaultBaseURL, secretName, secretVersion)
	if err != nil {
		log.Printf("Error getting secret: %v", err.Error())
		return Secret{}, err
	}

	result := Secret(secret)

	return result, nil
}

func GetSecretByURL(secretURL string) (Secret, error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		log.Printf("Failed to parse URL for secret: %v", err.Error())
		return Secret{}, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/secrets/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	secretName := res[0][1]

	result, err := GetSecret(vaultBaseURL, secretName, "")
	if err != nil {
		log.Printf("Failed to get secret from parsed values %v and %v: %v", vaultBaseURL, secretName, err.Error())
		return Secret{}, err
	}

	return result, nil
}

func GetSecrets(vaultBaseURL string) (results []Secret, err error) {
	max := int32(25)
	pages, err := newClient().GetSecrets(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting secret: %v", err.Error())
		return []Secret{}, err
	}

	for {
		for _, value := range pages.Values() {
			secretURL := *value.ID
			secret, err := GetSecretByURL(secretURL)
			if err != nil {
				log.Printf("Error loading secret contents: %v", err.Error())
				return nil, err
			}

			results = append(results, secret)
		}

		if pages.NotDone() {
			pages.Next()
		} else {
			break
		}
	}

	return results, nil
}
