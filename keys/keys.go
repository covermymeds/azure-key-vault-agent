package keys

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
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
