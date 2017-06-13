package main

import (
	"net/http"
	"MediaServer/urlgen"
	"MediaServer/websrv/station"
	"path"
	"log"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// someone wants to create a new station!
	
	// make a url and station
	newUrl := urlgen.Gen()
	liveStations[newUrl] = station.New(w, r)
	
	// dunno how "correct" this is
	http.Redirect(w, r, "/station/" + newUrl, http.StatusTemporaryRedirect)
}

func customHandler(path string, mimeType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := registry.Get(path)
		if data != nil {
			if mimeType != "" {
				w.Header().Add("Content-Type", mimeType)
			}
			
			w.Write(data)
		} else {
			log.Printf("Request for unknown file '%s'\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func handler(mimeType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := registry.Get(r.URL.Path)
		if data != nil {
			w.Header().Add("Content-Type", mimeType)
			w.Write(data)
		} else {
			log.Printf("Request for unknown file '%s'\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func songHandler(w http.ResponseWriter, r *http.Request) {
	// TODO request file from the database
	
	// tell client that it can request byte ranges of a song
	w.Header().Add("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)
}

func stationHandler(w http.ResponseWriter, r *http.Request) {
	stnUrl := path.Base(r.URL.Path)
	
	stn, ok := liveStations[stnUrl]
	if !ok {
		log.Printf("Request for non-existent station: '%s'\n", stnUrl)
		http.NotFound(w, r)
		return
	}
	
	tmpl := templatePages["/html/station.html"]
	
	data, err := tmpl.Generate(stn)
	if err != nil {
		log.Printf("Error generating html page for station '%s': '%s'", stnUrl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if data == nil || len(data) == 0 {
		log.Printf("No output generated for station page '%s'\n", stnUrl)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.Header().Add("Content-Type", "text/html")
	w.Write(data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO return results for specified search parameters
}
