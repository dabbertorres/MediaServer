package station

import (
	"bytes"
	"html/template"
	
	"MediaServer/internal/media"
)

type Page struct {
	template *template.Template
	length int
}

func Load(data []byte) Page {
	return Page{
		template: template.Must(template.New("stationPage").Parse(string(data))),
		length: len(data),
	}
}

func (p Page) Generate(station media.Station) ([]byte, error) {
	ret := make([]byte, 0, p.length)
	err := p.template.Execute(bytes.NewBuffer(ret), station)
	return ret, err
}
