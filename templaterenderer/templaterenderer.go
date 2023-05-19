package templaterenderer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/covermymeds/azure-key-vault-agent/certs"
	"github.com/covermymeds/azure-key-vault-agent/certutil"
	"github.com/covermymeds/azure-key-vault-agent/resource"
	"github.com/covermymeds/azure-key-vault-agent/secrets"
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
		"expandFullChain": func(items map[string]secrets.Secret) map[string]secrets.Secret {
			results := make(map[string]secrets.Secret)

			for secretName, secret := range items {
				results[secretName] = secret
				switch contentType := *secret.ContentType; contentType {
				case "application/x-pem-file":
					results[secretName+".key"] = cloneSecret(secret, certutil.PemPrivateKeyFromPem(*secret.Value))
					results[secretName+".pem"] = cloneSecret(secret, certutil.PemChainFromPem(*secret.Value, false))
				case "application/x-pkcs12":
					results[secretName+".key"] = cloneSecret(secret, certutil.PemPrivateKeyFromPkcs12(*secret.Value))
					results[secretName+".pem"] = cloneSecret(secret, certutil.PemChainFromPkcs12(*secret.Value, false))
				default:
					continue
				}
			}
			return results
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
		"toValues": func(secrets map[string]secrets.Secret) map[string]string {
			secretValues := make(map[string]string)
			for key, secret := range secrets {
				secretValues[key] = *secret.Value
			}
			return secretValues
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
		panic(fmt.Sprintf("Error executing template: %v Error: %v", templateContents, err))
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

func cloneSecret(secret secrets.Secret, parsedItem string) secrets.Secret {
	item := secret
	item.Value = nil
	item_copy := *secret.Value
	item.Value = &item_copy
	*item.Value = parsedItem

	return item
}
