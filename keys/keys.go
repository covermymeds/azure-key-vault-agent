package keys

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

func GetKey(vaultBaseURL string, keyName string, keyVersion string) (result keyvault.JSONWebKey, err error) {
	client := getClient()

	key, err := client.GetKey(context.Background(), vaultBaseURL, keyName, keyVersion)
	if err != nil {
		log.Fatalf("Error getting key: %v\n", err.Error())
	}

	result = *key.Key

	/*
		kb := *key.Key

		a := kb.Kty
		log.Println(a)
		b := kb.K
		log.Println(b)
		n := *kb.N
		log.Println(n)
		e := *kb.E
		log.Println(e)
		d := kb.D
		log.Println(d)
		dp := kb.DP
		log.Println(dp)
		dq := kb.DQ
		log.Println(dq)
		qi := kb.QI
		log.Println(qi)
		p := kb.P
		log.Println(p)
		q := kb.Q
		log.Println(q)
		k := kb.K
		log.Println(k)
		t := kb.T
		log.Println(t)

		result = ""
	*/

	return
}

func GetKeyByURL(keyURL string) (result keyvault.JSONWebKey, err error) {
	u, err := url.Parse(keyURL)
	if err != nil {
		log.Fatalf("Failed to parse URL for key: %v\n", err.Error())
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/keys/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	keyName := res[0][1]

	result, err = GetKey(vaultBaseURL, keyName, "")
	if err != nil {
		log.Fatalf("Failed to get key from parsed values %v and %v: %v\n", vaultBaseURL, keyName, err.Error())
	}

	return
}

func GetKeys(vaultBaseURL string) (results []keyvault.JSONWebKey, err error) {
	client := getClient()

	max := int32(25)
	pages, err := client.GetKeys(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Fatalf("Error getting key: %v\n", err.Error())
	}

	for {
		for _, value := range pages.Values() {
			keyURL := *value.Kid
			key, err := GetKeyByURL(keyURL)
			if err != nil {
				log.Fatalf("Error loading key contents: %v\n", err.Error())
			}

			results = append(results, key)
		}

		if pages.NotDone() {
			pages.Next()
		} else {
			break
		}
	}

	return
}
