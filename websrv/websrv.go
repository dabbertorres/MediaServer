package websrv

import (
	"path/filepath"

	"radio/cache"
)

// data needed at runtime to accomplish our desires!
var (
	registry     *cache.Registry
	liveStations = make(map[string]*Station)
)

// Init initializes data needed for running the web server
func Init(registryPath string) (err error) {
	registry, err = cache.NewRegistry(filepath.Clean(registryPath))
	if err != nil {
		return err
	}

	if err := registry.Walk(); err != nil {
		return err
	}

	helpTmpls := registry.Filter("tmpl/.*\\.tmpl")
	err = LoadHelperTemplates(helpTmpls...)
	if err != nil {
		return err
	}
	
	for _, f := range helpTmpls {
		go ReloadHelperTemplate(f.Name, registry.ListenTo(f.Name))
	}
	
	htmlTmpls := registry.Filter("html/.*\\.html")
	for _, f := range htmlTmpls {
		err = LoadTemplate(f)
		if err != nil {
			return err
		}
		
		go ReloadTemplate(f.Name, registry.ListenTo(f.Name))
	}

	return nil
}

// Closes releases resources held by the web server
func Close() {
	registry.Close()
}
