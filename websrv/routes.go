package websrv

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() http.Handler {
	router := mux.NewRouter()

	// application file paths

	router.Path("/").Methods("GET").HandlerFunc(indexHandler)
	router.PathPrefix("/css/").Methods("GET").HandlerFunc(handler("text/css"))
	router.PathPrefix("/html/").Methods("GET").HandlerFunc(handler("text/html"))
	router.PathPrefix("/js/").Methods("GET").HandlerFunc(handler("text/javascript"))
	router.PathPrefix("/img/png").Methods("GET").HandlerFunc(handler("image/png"))
	router.PathPrefix("/img/svg").Methods("GET").HandlerFunc(handler("image/svg+xml"))

	// favicons

	router.Path("/browserconfig.xml").Methods("GET").HandlerFunc(customHandler("favicon/browserconfig.xml", "application/xml"))
	router.Path("/manifest.json").Methods("GET").HandlerFunc(customHandler("favicon/manifest.json", "application/xml"))
	router.Path("/android-chrome-192x192.png").Methods("GET").HandlerFunc(customHandler("favicon/android-chrome-192x192.png", "image/png"))
	router.Path("/android-chrome-512x512.png").Methods("GET").HandlerFunc(customHandler("favicon/android-chrome-512x512.png", "image/png"))
	router.Path("/apple-touch-icon.png").Methods("GET").HandlerFunc(customHandler("favicon/apple-touch-icon.png", "image/png"))
	router.Path("/favicon.ico").Methods("GET").HandlerFunc(customHandler("favicon/favicon.ico", "image/x-icon"))
	router.Path("/favicon.png").Methods("GET").HandlerFunc(customHandler("favicon/favicon.png", "image/png"))
	router.Path("/favicon-16x16.png").Methods("GET").HandlerFunc(customHandler("favicon/favicon-16x16.png", "image/png"))
	router.Path("/favicon-32x32.png").Methods("GET").HandlerFunc(customHandler("favicon/favicon-32x32.png", "image/png"))
	router.Path("/mstile-150x150.png").Methods("GET").HandlerFunc(customHandler("favicon/mstile-150x150.png", "image/png"))
	router.Path("/safari-pinned-tab.svg").Methods("GET").HandlerFunc(customHandler("favicon/safari-pinned-tab.svg", "image/svg+xml"))

	// api paths

	search := router.PathPrefix("/search").Subrouter()
	search.Path("/").Methods("GET").HandlerFunc(searchHandler)

	song := router.PathPrefix("/song").Subrouter()
	song.Path("/{artist}/{title}").Methods("GET").HandlerFunc(songHandler)

	station := router.PathPrefix("/station").Subrouter()
	station.Path("/{name}").Methods("GET").HandlerFunc(tuneToStation)
	station.Path("/{name}").Methods("POST").HandlerFunc(startStation)
	station.Path("/{name}/playlist").Methods("GET").HandlerFunc(stationGetPlaylist)
	station.Path("/{name}/playlist").Methods("POST").HandlerFunc(stationAddSong)
	station.Path("/{name}/playlist").Methods("DELETE").HandlerFunc(stationRemoveSong)
	station.Path("/{name}/song").Methods("GET").HandlerFunc(stationGetPlayingSong)
	station.Path("/{name}/socket").Methods("GET").HandlerFunc(stationConnect)

	return router
}
