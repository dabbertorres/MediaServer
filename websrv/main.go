package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"radio/internal/interrupt"
	"radio/urlgen"
	"radio/websrv/file"
	"radio/websrv/station"
	"context"
)

const (
	appFileDir    = "app/"
	urlGenDataDir = "urlgen/dat"
)

var (
	registry *file.Registry

	liveStations = map[string]*station.Data{}

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
	
	waitShutdown := make(chan bool)

	// listen for termination signal
	interrupt.OnExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println("Server shutdown error:", err)
		}
		waitShutdown <- true
	})

	// hey now we can do what we want
	log.Println("Serving...")
	err = server.ListenAndServe()
	<-waitShutdown
	if err == http.ErrServerClosed {
		log.Println("Done")
	} else {
		log.Println("Server shutdown unexpectedly:", err)
	}
}
