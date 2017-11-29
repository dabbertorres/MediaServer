package websrv

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

		p, err := LoadPage(data)
		if err != nil {
			log.Printf("Error loading Page '%s': %v\n", path, err)
			continue
		}

		templatePages[path] = p
	}
}

type Page struct {
	template *template.Template
	length   int
}

func LoadPage(data []byte) (Page, error) {
	tmpl, err := template.New("page").Parse(string(data))
	return Page{
		template: tmpl,
		length:   len(data),
	}, err
}

func (p Page) Generate(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := p.template.Execute(&buf, data)
	return buf.Bytes(), err
}
