package templaterenderer

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/certutil"
	"github.com/covermymeds/azure-key-vault-agent/resource"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
	"io/ioutil"
	"text/template"
)

func RenderFile(path string, resourceMap resource.ResourceMap) string {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		panic(fmt.Sprintf("Error reading template %v: %v", path, err))
	}

	return RenderInline(string(contents), resourceMap)
}

func RenderInline(templateContents string, resourceMap resource.ResourceMap) string {
	helpers := template.FuncMap{
		"privateKey": func(secret secrets.Secret) string {
			switch contentType := *secret.ContentType; contentType {
			case "application/x-pem-file":
				return certutil.PemPrivateKeyFromPem(*secret.Value)
			case "application/x-pkcs12":
				return certutil.PemPrivateKeyFromPkcs12(*secret.Value)
			default:
				panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
			}
		},
		"cert": func(resource resource.Resource) string {
			switch t := resource.(type) {
			case certs.Cert:
				cert := resource.(certs.Cert)
				return certutil.PemCertFromBytes(*cert.Cer)
			case secrets.Secret:
				return certFromSecret(resource.(secrets.Secret))
			default:
				panic(fmt.Sprintf("Got unexpected type: %v", t))
			}
		},
		"issuers": func(secret secrets.Secret) string {
			switch contentType := *secret.ContentType; contentType {
			case "application/x-pem-file":
				return certutil.PemChainFromPem(*secret.Value, true)
			case "application/x-pkcs12":
				return certutil.PemChainFromPkcs12(*secret.Value, true)
			default:
				panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
			}
		},
		"fullChain": func(secret secrets.Secret) string {
			switch contentType := *secret.ContentType; contentType {
			case "application/x-pem-file":
				return certutil.PemChainFromPem(*secret.Value, false)
			case "application/x-pkcs12":
				return certutil.PemChainFromPkcs12(*secret.Value, false)
			default:
				panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
			}
		},
	}

	// Init the template
	t, err := template.New("template").Funcs(helpers).Funcs(sprig.TxtFuncMap()).Parse(templateContents)
	if err != nil {
		panic(fmt.Sprintf("Error parsing template:\n%v\nError:\n%v\n", templateContents, err))
	}

	// Execute the template
	var buf bytes.Buffer
	err = t.Execute(&buf, resourceMap)
	if err != nil {
		panic(fmt.Sprintf("Error executing template:\n%v\nResources:\n%v\nError:\n%v\n", templateContents, resourceMap, err))
	}

	result := buf.String()

	return result
}

func certFromSecret(secret secrets.Secret) string {
	switch contentType := *secret.ContentType; contentType {
	case "application/x-pem-file":
		return certutil.PemCertFromPem(*secret.Value)
	case "application/x-pkcs12":
		return certutil.PemCertFromPkcs12(*secret.Value)
	default:
		panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
	}
}