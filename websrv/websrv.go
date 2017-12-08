package websrv

import (
	"path/filepath"

	"radio/file"
)

// data needed at runtime to accomplish our desires!
var (
	registry     *file.Registry
	liveStations = make(map[string]*Station)
)

// Init initializes data needed for running the web server
func Init(registryPath string) error {
	reg, err := file.NewRegistry(filepath.Clean(registryPath))
	if err != nil {
		return err
	}

	if err := reg.Walk(); err != nil {
		return err
	}

	registry = reg

	loadTemplates()

	return nil
}

// Closes releases resources held by the web server
func Close() {
	registry.Close()
}
