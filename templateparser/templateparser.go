package templateparser

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"log"
	"path/filepath"
	"text/template"
)

func InlineTemplate(inline string, path string, resource resource.Resource) string {
	log.Printf("Parsing in-line template")
	// read in the template
	t := template.Must(template.New("inline").Funcs(sprig.TxtFuncMap()).Parse(inline))

	var buf bytes.Buffer
	// Execute the template
	err := t.Execute(&buf, resource)
	if err != nil {
		log.Panic(err)
	}

	result := buf.String()
	return result
}

func TemplateFile(tpath string, path string, resource resource.Resource) string {
	log.Printf("Parsing template file %v", tpath)
	// read in the template file
	t := template.Must(template.New(filepath.Base(tpath)).Funcs(sprig.TxtFuncMap()).ParseFiles(tpath))

	var buf bytes.Buffer
	// Execute the template
	err := t.Execute(&buf, resource)
	if err != nil {
		log.Panic(err)
	}

	result := buf.String()
	return result
}

