package certs

import (
	"context"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
)

type Cert keyvault.CertificateBundle

func (c Cert) String() string {
	return base64.StdEncoding.EncodeToString(*c.Cer)
}

func GetCert(client keyvault.BaseClient, vaultBaseURL string, certName string, certVersion string) (Cert, error) {
	cert, err := client.GetCertificate(context.Background(), vaultBaseURL, certName, certVersion)
	if err != nil {
		log.Printf("Error getting cert: %v", err.Error())
		return Cert{}, err
	}

	return Cert(cert), nil
}

func GetCertByURL(client keyvault.BaseClient, certURL string) (Cert, error) {
	u, err := url.Parse(certURL)
	if err != nil {
		log.Printf("Failed to parse URL for cert: %v", err.Error())
		return Cert{}, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/certificates/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	certName := res[0][1]

	result, err := GetCert(client, vaultBaseURL, certName, "")
	if err != nil {
		log.Printf("Failed to get cert from parsed values %v and %v: %v", vaultBaseURL, certName, err.Error())
		return Cert{}, err
	}

	return result, nil
}

func GetCerts(client keyvault.BaseClient, vaultBaseURL string) (results []Cert, err error) {
	max := int32(25)
	pages, err := client.GetCertificates(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting cert: %v", err.Error())
		return nil, err
	}

	for {
		for _, value := range pages.Values() {
			certURL := *value.ID
			cert, err := GetCertByURL(client, certURL)
			if err != nil {
				log.Printf("Error loading cert contents: %v", err.Error())
				return nil, err
			}

			results = append(results, cert)
		}

		if pages.NotDone() {
			pages.Next()
		} else {
			break
		}
	}

	return results, nil
}
