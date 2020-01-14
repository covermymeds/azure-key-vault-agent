package keys

import (
	"context"
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

type Key keyvault.KeyBundle

func (k Key) Map() map[string]interface{} {
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

func (k Key) String() string {
	//TODO
	return *k.Key.Kid
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

func GetKey(vaultBaseURL string, keyName string, keyVersion string) (resource.Resource, error) {
	key, err := newClient().GetKey(context.Background(), vaultBaseURL, keyName, keyVersion)
	if err != nil {
		log.Printf("Error getting key: %v\n", err.Error())
		return nil, err
	}

	result := *key.Key

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

	return result, err
}

func GetKeyByURL(keyURL string) (resource.Resource, error) {
	u, err := url.Parse(keyURL)
	if err != nil {
		log.Printf("Failed to parse URL for key: %v\n", err.Error())
		return nil, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/keys/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	keyName := res[0][1]

	result, err := GetKey(vaultBaseURL, keyName, "")
	if err != nil {
		log.Printf("Failed to get key from parsed values %v and %v: %v\n", vaultBaseURL, keyName, err.Error())
		return nil, err
	}

	return result, nil
}

func GetKeys(vaultBaseURL string) ([]resource.Resource, error) {
	max := int32(25)
	pages, err := newClient().GetKeys(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting key: %v\n", err.Error())
		return nil, err
	}

	var results []resource.Resource
	for {
		for _, value := range pages.Values() {
			keyURL := *value.Kid
			key, err := GetKeyByURL(keyURL)
			if err != nil {
				log.Printf("Error loading key contents: %v\n", err.Error())
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
