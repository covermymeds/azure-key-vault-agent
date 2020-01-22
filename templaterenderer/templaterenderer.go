package templaterenderer

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/chrisjohnson/azure-key-vault-agent/certutil"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"io/ioutil"
	"log"
	"text/template"
)

func RenderFile(path string, resourceMap resource.ResourceMap) string {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		log.Panicf("Error reading template %v: %v", path, err)
	}

	return RenderInline(string(contents), resourceMap)
}

func RenderInline(templateContents string, resourceMap resource.ResourceMap) string {
	helpers := template.FuncMap{
		"privateKey": func(name string) interface{} {
			value, ok := resourceMap.Secrets[name]
			privateKey := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					privateKey = certutil.PemPrivateKeyFromPem(*value.Value)
				case "application/x-pkcs12":
					privateKey = certutil.PemPrivateKeyFromPkcs12(*value.Value)
				default:
					log.Panicf("Got unexpected content type: %v", contentType)
				}
			} else {
				log.Panicf("privateKey lookup failed: Expected a Secret with name %v\n", name)
			}
			return privateKey
		},
		"cert": func(name string) interface{} {
			// TODO: If the cert can be found on either a Cert or a Secret, we should handle discovering it from both
			value, ok := resourceMap.Secrets[name]
			cert := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					cert = certutil.PemCertFromPem(*value.Value)
				case "application/x-pkcs12":
					cert = certutil.PemCertFromPkcs12(*value.Value)
				default:
					log.Panicf("Got unexpected content type: %v", contentType)
				}
			} else {
				log.Panicf("cert lookup failed: Expected a Secret with name %v\n", name)
			}
			return cert
		},
		"chain": func(name string) interface{} {
			value, ok := resourceMap.Secrets[name]
			chain := ""
			if ok {
				switch contentType := *value.ContentType; contentType {
				case "application/x-pem-file":
					chain = certutil.PemChainFromPem(*value.Value)
				case "application/x-pkcs12":
					chain = certutil.PemChainFromPkcs12(*value.Value)
				default:
					log.Panicf("Got unexpected content type: %v", contentType)
				}
			} else {
				log.Panicf("cert lookup failed: Expected a Secret with name %v\n", name)
			}
			return chain
		},
	}

	// Init the template
	t, err := template.New("template").Funcs(helpers).Funcs(sprig.TxtFuncMap()).Parse(templateContents)
	if err != nil {
		log.Panicf("Error parsing template:\n%v\nError:\n%v\n", templateContents, err)
	}

	// Execute the template
	var buf bytes.Buffer
	err = t.Execute(&buf, resourceMap)
	if err != nil {
		log.Panicf("Error executing template:\n%v\nResources:\n%v\nError:\n%v\n", templateContents, resourceMap, err)
	}

	result := buf.String()

	return result
}

