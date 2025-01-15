package client

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
	"github.com/covermymeds/azure-key-vault-agent/keys"
	"github.com/covermymeds/azure-key-vault-agent/config"
	"github.com/covermymeds/azure-key-vault-agent/iam"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
)

type KeyvaultClient struct {
	Client keyvault.BaseClient
}

func NewKeyvaultClient(cred config.CredentialConfig) KeyvaultClient {
	c := keyvault.New()
	authorizer, err := iam.GetKeyvaultAuthorizer(cred.TenantID, cred.ClientID, cred.ClientSecret)
	if err != nil {
		panic(fmt.Sprintf("Error authorizing: %v", err.Error()))
	}
	c.Authorizer = authorizer
	return KeyvaultClient{c}
}

func (c KeyvaultClient) GetCert(vaultBaseURL string, certName string, certVersion string) (certs.Cert, error) {
	cert, err := c.Client.GetCertificate(context.Background(), vaultBaseURL, certName, certVersion)
	if err != nil {
		log.Printf("Error getting cert: %v", err.Error())
		return certs.Cert{}, err
	}

	return certs.Cert(cert), nil
}

func (c KeyvaultClient) GetCertByURL(certURL string) (certs.Cert, error) {
	u, err := url.Parse(certURL)
	if err != nil {
		log.Printf("Failed to parse URL for cert: %v", err.Error())
		return certs.Cert{}, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/certificates/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	certName := res[0][1]

	result, err := c.GetCert(vaultBaseURL, certName, "")
	if err != nil {
		log.Printf("Failed to get cert from parsed values %v and %v: %v", vaultBaseURL, certName, err.Error())
		return certs.Cert{}, err
	}

	return result, nil
}

func (c KeyvaultClient) GetCerts(vaultBaseURL string) (results []certs.Cert, err error) {
	max := int32(25)
	pages, err := c.Client.GetCertificates(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting cert: %v", err.Error())
		return nil, err
	}

	for {
		for _, value := range pages.Values() {
			certURL := *value.ID
			cert, err := c.GetCertByURL(certURL)
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

func (c KeyvaultClient) GetSecret(vaultBaseURL string, secretName string, secretVersion string) (secrets.Secret, error) {
	secret, err := c.Client.GetSecret(context.Background(), vaultBaseURL, secretName, secretVersion)
	if err != nil {
		log.Printf("Error getting secret: %v", err.Error())
		return secrets.Secret{}, err
	}

	result := secrets.Secret{
		Value: secret.Value,
		ContentType: secret.ContentType,
	}

	return result, nil
}

func (c KeyvaultClient) GetSecretByURL(secretURL string) (string, secrets.Secret, error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		log.Printf("Failed to parse URL for secret: %v", err.Error())
		return "invalid", secrets.Secret{}, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/secrets/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	secretName := res[0][1]

	result, err := c.GetSecret(vaultBaseURL, secretName, "")
	if err != nil {
		log.Printf("Failed to get secret from parsed values %v and %v: %v", vaultBaseURL, secretName, err.Error())
		return "invalid", secrets.Secret{}, err
	}

	return secretName, result, nil
}

func (c KeyvaultClient) GetSecrets(vaultBaseURL string) (results map[string]secrets.Secret, err error) {
	max := int32(25)
	pages, err := c.Client.GetSecrets(context.Background(), vaultBaseURL, &max)
	results = make(map[string]secrets.Secret)
	if err != nil {
		log.Printf("Error getting secret: %v", err.Error())
		return map[string]secrets.Secret{}, err
	}

	for {
		for _, value := range pages.Values() {
			if *value.Attributes.Enabled {
				secretURL := *value.ID
				secretName, secret, err := c.GetSecretByURL(secretURL)

				if err != nil {
					log.Printf("Error loading secret contents: %v", err.Error())
					return nil, err
				}

				results[secretName] = secret
			}

		}

		if pages.NotDone() {
			pages.Next()
		} else {
			break
		}
	}

	return results, nil
}

func (c KeyvaultClient) GetKey(vaultBaseURL string, keyName string, keyVersion string) (keys.Key, error) {
	key, err := c.Client.GetKey(context.Background(), vaultBaseURL, keyName, keyVersion)
	if err != nil {
		log.Printf("Error getting key: %v", err.Error())
		return keys.Key{}, err
	}

	result := keys.Key(key)

	return result, err
}

func (c KeyvaultClient) GetKeyByURL(keyURL string) (keys.Key, error) {
	u, err := url.Parse(keyURL)
	if err != nil {
		log.Printf("Failed to parse URL for key: %v", err.Error())
		return keys.Key{}, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/keys/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	keyName := res[0][1]

	result, err := c.GetKey(vaultBaseURL, keyName, "")
	if err != nil {
		log.Printf("Failed to get key from parsed values %v and %v: %v", vaultBaseURL, keyName, err.Error())
		return keys.Key{}, err
	}

	return result, nil
}

func (c KeyvaultClient) GetKeys(vaultBaseURL string) (results []keys.Key, err error) {
	max := int32(25)
	pages, err := c.Client.GetKeys(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting key: %v", err.Error())
		return nil, err
	}

	for {
		for _, value := range pages.Values() {
			keyURL := *value.Kid
			key, err := c.GetKeyByURL(keyURL)
			if err != nil {
				log.Printf("Error loading key contents: %v", err.Error())
				return nil, err
			}

			results = append(results, key)
		}

		if pages.NotDone() {
			pages.Next()
		} else {
			break
		}
	}

	return results, nil
}
