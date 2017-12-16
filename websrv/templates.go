package websrv

import (
	"bytes"
	"html/template"
	"log"
	"sync"

	"radio/cache"
	"errors"
)

var (
	htmlTemplates  = make(map[string]*template.Template)
	helperTemplates = make(map[string]*template.Template)
	templatesMutex sync.RWMutex
	
	TemplateDoesNotExist = errors.New("template does not exist")
)

func LoadTemplate(file cache.File) error {
	t, err := template.New(file.Name).Parse(string(file.Data))
	if err != nil {
		return err
	}

	for _, ht := range helperTemplates {
		t, err = t.AddParseTree(ht.Name(), ht.Tree)
		if err != nil {
			return err
		}
	}

	htmlTemplates[file.Name] = t

	return nil
}

func LoadHelperTemplates(files ...cache.File) error {
	for _, f := range files {
		t, err := template.New(f.Name).Parse(string(f.Data))
		if err != nil {
			return err
		}
		
		helperTemplates[f.Name] = t
	}

	return nil
}

func ReloadTemplate(name string, c <-chan []byte) {
	var err error
	
	for d := range c {
		templatesMutex.Lock()
		htmlTemplates[name], err = htmlTemplates[name].Parse(string(d))
		templatesMutex.Unlock()
		
		if err != nil {
			log.Printf("Error reloading helper template '%s': %v\n", name, err)
		}
	}
}

func ReloadHelperTemplate(name string, c <-chan []byte) {
	var err error
	
	for d := range c {
		templatesMutex.Lock()
		helperTemplates[name], err = helperTemplates[name].Parse(string(d))
		templatesMutex.Unlock()
		
		if err != nil {
			log.Printf("Error reloading helper template '%s': %v\n", name, err)
		}
	}
}

func RunTemplate(name string, data interface{}) ([]byte, error) {
	var (
		err error
		buf bytes.Buffer
	)
	
	templatesMutex.RLock()
	if t, ok := htmlTemplates[name]; ok {
		err = t.ExecuteTemplate(&buf, name, data)
	} else {
		err = TemplateDoesNotExist
	}
	templatesMutex.RUnlock()
	
	return buf.Bytes(), err
}
