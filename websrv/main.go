package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"MediaServer/internal/interrupt"
	"MediaServer/urlgen"
	"MediaServer/websrv/file"
	"MediaServer/websrv/station"
)

const (
	appFileDir    = "app/"
	urlGenDataDir = "urlgen/dat"
)

var (
	registry *file.Registry

	liveStations = map[string]station.Data{}

	// annoying to type, short for "sanitize"
	san = filepath.Clean
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	var err error

	registry, err = file.NewRegistry(appFileDir)
	if err != nil {
		panic(err)
	}
	defer registry.Close()

	err = registry.Walk(nil)
	if err != nil {
		panic(err)
	}

	err = urlgen.LoadDir(san(urlGenDataDir))
	if err != nil {
		panic(err)
	}
	
	loadTemplates()

	server := http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 32,
		Handler:        makeMuxer(),
	}

	// listen for termination signal!
	interrupt.OnExit(func() { server.Close() })

	// hey now we can do what we want
	fmt.Println("Serving...")
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Done")
	} else {
		log.Println("Server shutdown unexpectedly:", err)
	}
}
