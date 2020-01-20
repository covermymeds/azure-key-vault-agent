package templateparser

import (
	"bytes"
	"log"
	"path/filepath"
	"text/template"

	"github.com/chrisjohnson/azure-key-vault-agent/resource"

	"github.com/Masterminds/sprig"
)

func InlineTemplate(inline string, path string, resource resource.ResourceMap) string {
	// Read in the template
	t, err := template.New("inline").Funcs(sprig.TxtFuncMap()).Parse(inline)
	if err != nil {
		log.Fatalf("Error parsing template:\n%v\nError:\n%v\n", inline, err)
	}

	var buf bytes.Buffer
	// Execute the template
	err = t.Execute(&buf, resource)
	if err != nil {
		log.Fatalf("Error executing template:\n%v\nResource:\n%v\nError:\n%v\n", inline, resource, err)
	}

	result := buf.String()
	return result
}

func TemplateFile(tpath string, path string, resource resource.ResourceMap) string {
	// Read in the template file
	t, err := template.New(filepath.Base(tpath)).Funcs(sprig.TxtFuncMap()).ParseFiles(tpath)
	if err != nil {
		log.Fatalf("Error parsing template:\n%v\nError:\n%v\n", tpath, err)
	}

	var buf bytes.Buffer
	// Execute the template
	err = t.Execute(&buf, resource)
	if err != nil {
		log.Fatalf("Error executing template:\n%v\nResource:\n%v\nError:\n%v\n", tpath, resource, err)
	}

	result := buf.String()
	return result
}
