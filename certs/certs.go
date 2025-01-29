package certs

import (
	"encoding/base64"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
)

type Cert keyvault.CertificateBundle

func (c Cert) String() string {
	return base64.StdEncoding.EncodeToString(*c.Cer)
}
