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
	"log"
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
		"privateKey": func(name string) string {
			value, ok := resourceMap.Secrets[name]
			privateKey := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					privateKey = certutil.PemPrivateKeyFromPem(*value.Value)
				case "application/x-pkcs12":
					privateKey = certutil.PemPrivateKeyFromPkcs12(*value.Value)
				default:
					panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
				}
			} else {
				panic(fmt.Sprintf("privateKey lookup failed: Expected a Secret with name %v", name))
			}
			return privateKey
		},
		"cert": func(resource resource.Resource) string {
			switch t := resource.(type) {
			case certs.Cert:
				cert := resource.(certs.Cert)
				return certutil.PemCertFromBytes(*cert.Cer)
			case secrets.Secret:
				log.Print("Got a Secret")
				return certFromSecret(resource.(secrets.Secret))
			default:
				panic(fmt.Sprintf("Got unexpected type: %v", t))
			}
		},
		"issuers": func(name string) string {
			value, ok := resourceMap.Secrets[name]
			issuers := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					issuers = certutil.PemChainFromPem(*value.Value, true)
				case "application/x-pkcs12":
					issuers = certutil.PemChainFromPkcs12(*value.Value, true)
				default:
					panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
				}
			} else {
				panic(fmt.Sprintf("cert lookup failed: Expected a Secret with name %v", name))
			}
			return issuers
		},
		"fullChain": func(name string) string {
			value, ok := resourceMap.Secrets[name]
			fullChain := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					fullChain = certutil.PemChainFromPem(*value.Value, false)
				case "application/x-pkcs12":
					fullChain = certutil.PemChainFromPkcs12(*value.Value, false)
				default:
					panic(fmt.Sprintf("Got unexpected content type: %v", contentType))
				}
			} else {
				panic(fmt.Sprintf("cert lookup failed: Expected a Secret with name %v", name))
			}
			return fullChain
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
