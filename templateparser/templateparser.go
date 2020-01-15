package templateparser

import (
	"github.com/Masterminds/sprig"
	"github.com/chrisjohnson/azure-key-vault-agent/resource"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func InlineTemplate(inline string, path string, resource resource.Resource) {
	// read in the template file
	t := template.Must(template.New("inline").Funcs(sprig.TxtFuncMap()).Parse(inline))

	// create the destination file
	f, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	// Execute the template
	err = t.Execute(f, resource)
	if err != nil {
		log.Panic(err)
	}
}

func TemplateFile(tpath string, path string, resource resource.Resource) {
	// read in the template file
	t := template.Must(template.New(filepath.Base(tpath)).Funcs(sprig.TxtFuncMap()).ParseFiles(tpath))

	// create the destination file
	f, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	// Execute the template
	err = t.Execute(f, resource)
	if err != nil {
		log.Panic(err)
	}
}

