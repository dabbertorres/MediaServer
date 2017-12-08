package websrv

import (
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"

	"radio/urlgen"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
	// TODO make index.html
	// should have a small description, links, etc
	// have login, register, new temp station choices

	// user accounts (optional) basically reserve a permanent, and custom station name

	// make a url and station
	newUrl := urlgen.Gen()
	liveStations[newUrl] = NewStation(newUrl)

	http.Redirect(w, r, "/station/"+newUrl, http.StatusTemporaryRedirect)
}

func customHandler(path string, mimeType string) http.HandlerFunc {
	path = filepath.Clean(path)
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

func handler(mimeType string) http.HandlerFunc {
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
	panic("TODO")
	// TODO request file from the database

	// tell client that it can request byte ranges of a song
	w.Header().Add("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)
}

func tuneToStation(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	stn, ok := liveStations[name]
	if !ok {
		log.Printf("Request for non-existent station: '%s'\n", name)
		http.NotFound(w, r)
		return
	}

	tmpl := templatePages["/html/station.html"]

	data, err := tmpl.Generate(stn)
	if err != nil {
		log.Printf("Error generating html page for station '%s': '%s'", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if data == nil || len(data) == 0 {
		log.Printf("No output generated for station page '%s'\n", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write(data)
}

func startStation(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func stationGetPlaylist(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func stationAddSong(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func stationRemoveSong(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func stationGetPlayingSong(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO return results for specified search parameters
	panic("TODO")
}

func stationConnect(w http.ResponseWriter, r *http.Request) {
	stnUrl := path.Base(r.URL.Path)

	stn, ok := liveStations[stnUrl]
	if !ok {
		log.Printf("Socket request for non-existent station: '%s'\n", stnUrl)
		http.NotFound(w, r)
		return
	}

	stn.Add(w, r)
}
