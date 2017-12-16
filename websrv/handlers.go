package websrv

import (
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"

	"radio/urlgen"
)

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

func templateHandler(path, mimeType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO need to pass stuff to the template
		data, err := RunTemplate(path, nil)
		if err != nil {
			if err == TemplateDoesNotExist {
				log.Printf("Request for unknown file '%s'\n", path)
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Printf("Error executing template '%s': %v\n", path, err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			
			return
		}

		w.Header().Add("Content-Type", mimeType)
		w.Write(data)
	}
}

func customHandler(path, mimeType string) http.HandlerFunc {
	path = filepath.Clean(path)
	return func(w http.ResponseWriter, r *http.Request) {
		data := registry.Get(path)
		if data != nil {
			w.Header().Add("Content-Type", mimeType)
			w.Write(data)
		} else {
			log.Printf("Request for unknown '%s'\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func songHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
	// TODO request song file from the database

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

	// TODO upgrade to websocket connection, get user info, etc
	stn.TuneIn <- NewClient("", nil)

	// TODO data to pass!
	data, err := RunTemplate("html/station.html", nil)
	if err != nil {
		log.Printf("Error generating html page for station '%s': '%s'", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data == nil || len(data) == 0 {
		log.Printf("No output generated for station page '%s'\n", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write(data)
}

func startStation(w http.ResponseWriter, r *http.Request) {
	// TODO check if session corresponds to a user
	// if so, then don't generate a new station
	// just start up that user's station

	// make a url and station
	newUrl := urlgen.Gen()
	liveStations[newUrl] = NewStation(newUrl)

	http.Redirect(w, r, "/station/"+newUrl, http.StatusTemporaryRedirect)
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
