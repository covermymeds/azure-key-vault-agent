package certs

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	//"reflect"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"

	"github.com/chrisjohnson/azure-key-vault-agent/config"
	"github.com/chrisjohnson/azure-key-vault-agent/iam"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
)

type Cert keyvault.CertificateBundle

func (c Cert) Map() map[string]interface{} {
	m := make(map[string]interface{})
	/*
		v := reflect.ValueOf(m)
		for i := 0; i < v.NumField(); i++ {
			log.Println(v.Field(i))
			//m[v.Field(i)] = v.Field(i).Interface()
		}
	*/

	return m
}

func (c Cert) String() string {
	return base64.StdEncoding.EncodeToString(*c.Cer)
}

func newClient() keyvault.BaseClient {
	client := keyvault.New()
	a, err := iam.GetKeyvaultAuthorizer()
	if err != nil {
		log.Panicf("Error authorizing: %v\n", err.Error())
	}
	client.Authorizer = a
	client.AddToUserAgent(config.UserAgent())
	return client
}

func GetCert(vaultBaseURL string, certName string, certVersion string) (resource.Resource, error) {
	cert, err := newClient().GetCertificate(context.Background(), vaultBaseURL, certName, certVersion)
	if err != nil {
		log.Printf("Error getting cert: %v\n", err.Error())
		return nil, err
	}

	return Cert(cert), nil
}

func GetCertByURL(certURL string) (resource.Resource, error) {
	u, err := url.Parse(certURL)
	if err != nil {
		log.Printf("Failed to parse URL for cert: %v\n", err.Error())
		return nil, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/certificates/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	certName := res[0][1]

	result, err := GetCert(vaultBaseURL, certName, "")
	if err != nil {
		log.Printf("Failed to get cert from parsed values %v and %v: %v\n", vaultBaseURL, certName, err.Error())
		return nil, err
	}

	return result, nil
}

func GetCerts(vaultBaseURL string) ([]resource.Resource, error) {
	max := int32(25)
	pages, err := newClient().GetCertificates(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting cert: %v\n", err.Error())
		return nil, err
	}

	var results []resource.Resource
	for {
		for _, value := range pages.Values() {
			certURL := *value.ID
			cert, err := GetCertByURL(certURL)
			if err != nil {
				log.Printf("Error loading cert contents: %v\n", err.Error())
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
