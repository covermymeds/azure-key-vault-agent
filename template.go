package main

import (
	"log"
	"os"
	"text/template"
	"github.com/Masterminds/sprig"
)

func main() {

	fmap := sprig.TxtFuncMap()
	tpath := "test.tmpl"

	// read in the template file
	t := template.Must(template.New(tpath).Funcs(fmap).ParseFiles(tpath))


	// create the destination file
	path := "rendered.txt"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Prepare some data to insert into the template.
	type Recipient struct {
		Name, Gift string
		Attended   bool
	}
	var recipients = []Recipient{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute(f, r)
		if err != nil {
			log.Fatal(err)
		}
	}

}