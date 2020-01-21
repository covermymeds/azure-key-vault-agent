package templaterenderer

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/chrisjohnson/azure-key-vault-agent/resource"

	"github.com/Masterminds/sprig"
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
			if ok {
				// TODO: Transform value to extract the private key using some sort of library that can parse PEM format
				// TODO: How to handle PKCS12?
			} else {
				log.Panicf("privateKey lookup failed: Expected a Secret with name %v\n", name)
			}
			return value
		},
		"cert": func(name string) interface{} {
			value, ok := resourceMap.Secrets[name]
			// TODO: If the cert can be found on either a Cert or a Secret, we should handle discovering it from both
			if ok {
				// TODO: Transform value to extract the cert using some sort of library that can parse PEM format
				// TODO: How to handle PKCS12?
			} else {
				log.Panicf("cert lookup failed: Expected a Secret with name %v\n", name)
			}
			return value
		},
	}

	// Read in the template
	t, err := template.New("template").Funcs(helpers).Funcs(sprig.TxtFuncMap()).Parse(templateContents)
	if err != nil {
		log.Panicf("Error parsing template:\n%v\nError:\n%v\n", templateContents, err)
	}

	var buf bytes.Buffer
	// Execute the template
	err = t.Execute(&buf, resourceMap)
	if err != nil {
		log.Panicf("Error executing template:\n%v\nResource:\n%v\nError:\n%v\n", templateContents, resourceMap, err)
	}

	result := buf.String()

	return result
}
