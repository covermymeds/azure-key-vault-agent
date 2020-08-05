package keys

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"

	"github.com/chrisjohnson/azure-key-vault-agent/iam"
)

type Key keyvault.KeyBundle

func (k Key) String() string {
	bytes, _ := k.MarshalJSON()
	return string(bytes)
}

// MarshalJSON is the custom marshaler for KeyBundle.
func (kb Key) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if kb.Key != nil {
		objectMap["key"] = kb.Key
	}
	if kb.Attributes != nil {
		objectMap["attributes"] = kb.Attributes
	}
	if kb.Tags != nil {
		objectMap["tags"] = kb.Tags
	}
	return json.Marshal(objectMap)
}

func newClient() keyvault.BaseClient {
	client := keyvault.New()
	a, err := iam.GetKeyvaultAuthorizer()
	if err != nil {
		panic(fmt.Sprintf("Error authorizing: %v", err.Error()))
	}
	client.Authorizer = a
	return client
}

func GetKey(vaultBaseURL string, keyName string, keyVersion string) (Key, error) {
	key, err := newClient().GetKey(context.Background(), vaultBaseURL, keyName, keyVersion)
	if err != nil {
		log.Printf("Error getting key: %v", err.Error())
		return Key{}, err
	}

	result := Key(key)

	return result, err
}

func GetKeyByURL(keyURL string) (Key, error) {
	u, err := url.Parse(keyURL)
	if err != nil {
		log.Printf("Failed to parse URL for key: %v", err.Error())
		return Key{}, err
	}
	vaultBaseURL := fmt.Sprintf("%v://%v", u.Scheme, u.Host)

	regex := *regexp.MustCompile(`/keys/(.*)(/.*)?`)
	res := regex.FindAllStringSubmatch(u.Path, -1)
	keyName := res[0][1]

	result, err := GetKey(vaultBaseURL, keyName, "")
	if err != nil {
		log.Printf("Failed to get key from parsed values %v and %v: %v", vaultBaseURL, keyName, err.Error())
		return Key{}, err
	}

	return result, nil
}

func GetKeys(vaultBaseURL string) (results []Key, err error) {
	max := int32(25)
	pages, err := newClient().GetKeys(context.Background(), vaultBaseURL, &max)
	if err != nil {
		log.Printf("Error getting key: %v", err.Error())
		return nil, err
	}

	for {
		for _, value := range pages.Values() {
			keyURL := *value.Kid
			key, err := GetKeyByURL(keyURL)
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
