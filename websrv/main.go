package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"MediaServer/websrv/file"
	"MediaServer/websrv/interrupt"
)

type FileData []byte
type FileDataTree map[string]*FileData

// TODO
type Station struct {
}

var (
	registry *file.Registry

	// annoying to type, short for "sanitize"
	san = filepath.Clean
)

//var stations = make(map[string]Station)

func main() {
	var err error

	registry, err = file.NewRegistry("app")
	if err != nil {
		panic(err)
	}
	defer registry.Close()

	err = registry.Walk(nil)
	if err != nil {
		panic(err)
	}

	serverMux := http.NewServeMux()

	// basic server files!
	serverMux.HandleFunc("/", customHandler(san("app/html/index.html"), "text/html"))
	serverMux.HandleFunc("/css/", handler("text/css"))
	serverMux.HandleFunc("/html/", handler("text/html"))
	serverMux.HandleFunc("/js/", handler("text/javascript"))

	// images!
	serverMux.HandleFunc("/img/png", handler("image/png"))
	serverMux.HandleFunc("/img/svg", handler("image/svg+xml"))

	// favicon config crap (bless you, realfavicongenerator.net)
	serverMux.HandleFunc("/browserconfig.xml", customHandler(san("app/browserconfig.xml"), "application/xml"))
	serverMux.HandleFunc("/manifest.json", customHandler(san("app/manifest.json"), "application/json"))

	// actual favicon s
	serverMux.HandleFunc("/android-chrome-192x192.png", customHandler(san("app/img/favicon/android-chrome-192x192.png"), "image/png"))
	serverMux.HandleFunc("/android-chrome-512x512.png", customHandler(san("app/img/favicon/android-chrome-512x512.png"), "image/png"))
	serverMux.HandleFunc("/apple-touch-icon.png", customHandler(san("app/img/favicon/apple-touch-icon.png"), "image/png"))
	serverMux.HandleFunc("/favicon.ico", customHandler(san("app/img/favicon/favicon.ico"), "image/x-icon"))
	serverMux.HandleFunc("/favicon.png", customHandler(san("app/img/favicon/favicon.png"), "image/png"))
	serverMux.HandleFunc("/favicon-16x16.png", customHandler(san("app/img/favicon/favicon-16x16.png"), "image/png"))
	serverMux.HandleFunc("/favicon-32x32.png", customHandler(san("app/img/favicon/favicon-32x32.png"), "image/png"))
	serverMux.HandleFunc("/mstile-150x150.png", customHandler(san("app/img/favicon/mstile-150x150.png"), "image/png"))
	serverMux.HandleFunc("/safari-pinned-tab.svg", customHandler(san("app/img/favicon/safari-pinned-tab.svg"), "image/svg"))

	// actually interesting stuff eventually
	serverMux.HandleFunc("/media/", mediaHandler)
	serverMux.HandleFunc("/station/", stationHandler)

	server := http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 32,
		Handler:        serverMux,
	}

	// listen for termination signal!
	interrupt.OnExit(func() { server.Close() })

	// hey now we can do what we want
	fmt.Println("Serving...")
	server.ListenAndServe()
	fmt.Println("Done")
}

func customHandler(path string, mimeType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		data := registry.Get(path)
		if data != nil {
			if mimeType != "" {
				w.Header().Add("Content-Type", mimeType)
			}

			w.Write(data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func handler(mimeType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(registry.BasePath, r.URL.EscapedPath())
		path = filepath.Clean(path)

		data := registry.Get(path)
		if data != nil {
			w.Header().Add("Content-Type", mimeType)
			w.Write(data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	// TODO request file from the database (POST)
}

func stationHandler(w http.ResponseWriter, r *http.Request) {
	// TODO return specified station (GET)
}
