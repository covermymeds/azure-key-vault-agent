package certs

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

func GetCert(vaultBaseURL string, certName string, certVersion string) (result string, err error) {
	client := getClient()

	cert, err := client.GetCertificate(context.Background(), vaultBaseURL, certName, certVersion)
	if err != nil {
		log.Fatalf("Error getting cert: %v\n", err.Error())
	}

	// TODO: Return bundle?
	result = *cert.X509Thumbprint

	return
}

func GetCertByURL(certURL string) (result string, err error) {
	u, err := url.Parse(certURL)
	if err != nil {
		log.Fatalf("Failed to parse URL for cert: %v\n", err.Error())
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/certificates/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	certName := res[0][1]

	result, err = GetCert(vaultBaseURL, certName, "")
	if err != nil {
		log.Fatalf("Failed to get cert from parsed values %v and %v: %v\n", vaultBaseURL, certName, err.Error())
	}

	return
}

func GetCerts(vaultBaseURL string) (results []string, err error) {
	client := getClient()

	max := int32(25)
	pages, err := client.GetCertificates(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Fatalf("Error getting cert: %v\n", err.Error())
	}

	for {
		for _, value := range pages.Values() {
			certURL := *value.ID
			cert, err := GetCertByURL(certURL)
			if err != nil {
				log.Fatalf("Error loading cert contents: %v\n", err.Error())
			}

			results = append(results, cert)
		}

		if pages.NotDone() {
			pages.Next()
		} else {
			break
		}
	}

	return
}
