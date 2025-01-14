package secrets

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
)

type Secret keyvault.SecretBundle

func (s Secret) String() string {
	return *s.Value
}
