package main

import (
	"bytes"
	"html/template"
	"log"
)

var templatePages = map[string]Page{
	"/html/station.html": {},
}

func loadTemplates() {
	for path := range templatePages {
		data := registry.Get(path)
		if data == nil {
			log.Printf("Template page file '%s' not found!\n", path)
			continue
		}

		templatePages[path] = LoadPage(data)
	}
}

type Page struct {
	template *template.Template
	length   int
}

func LoadPage(data []byte) Page {
	return Page{
		template: template.Must(template.New("page").Parse(string(data))),
		length:   len(data),
	}
}

func (p Page) Generate(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := p.template.Execute(&buf, data)
	return buf.Bytes(), err
}
