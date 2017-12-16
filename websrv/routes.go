package websrv

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	MimeTypeHTML       = "text/html"
	MimeTypeCSS        = "text/css"
	MimeTypeJavaScript = "text/javascript"
	MimeTypePNG        = "image/png"
	MimeTypeSVG        = "image/svg+xml"
	MimeTypeXML        = "application/xml"
	MimeTypeICO        = "image/x-icon"
	MimeTypeJSON       = "application/json"
)

func Routes() http.Handler {
	router := mux.NewRouter()

	// application cache paths

	router.Path("/").Methods("GET").HandlerFunc(templateHandler("html/index.html", MimeTypeHTML))
	router.PathPrefix("/css/").Methods("GET").HandlerFunc(handler(MimeTypeCSS))
	router.PathPrefix("/html/").Methods("GET").HandlerFunc(handler(MimeTypeHTML))
	router.PathPrefix("/js/").Methods("GET").HandlerFunc(handler(MimeTypeJavaScript))
	router.PathPrefix("/img/png/").Methods("GET").HandlerFunc(handler(MimeTypePNG))
	router.PathPrefix("/img/svg/").Methods("GET").HandlerFunc(handler(MimeTypeSVG))

	// favicons

	router.Path("/browserconfig.xml").Methods("GET").HandlerFunc(customHandler("favicon/browserconfig.xml", MimeTypeXML))
	router.Path("/manifest.json").Methods("GET").HandlerFunc(customHandler("favicon/manifest.json", MimeTypeJSON))
	router.Path("/android-chrome-192x192.png").Methods("GET").HandlerFunc(customHandler("favicon/android-chrome-192x192.png", MimeTypePNG))
	router.Path("/android-chrome-512x512.png").Methods("GET").HandlerFunc(customHandler("favicon/android-chrome-512x512.png", MimeTypePNG))
	router.Path("/apple-touch-icon.png").Methods("GET").HandlerFunc(customHandler("favicon/apple-touch-icon.png", MimeTypePNG))
	router.Path("/favicon.ico").Methods("GET").HandlerFunc(customHandler("favicon/favicon.ico", MimeTypeICO))
	router.Path("/favicon.png").Methods("GET").HandlerFunc(customHandler("favicon/favicon.png", MimeTypePNG))
	router.Path("/favicon-16x16.png").Methods("GET").HandlerFunc(customHandler("favicon/favicon-16x16.png", MimeTypePNG))
	router.Path("/favicon-32x32.png").Methods("GET").HandlerFunc(customHandler("favicon/favicon-32x32.png", MimeTypePNG))
	router.Path("/mstile-150x150.png").Methods("GET").HandlerFunc(customHandler("favicon/mstile-150x150.png", MimeTypePNG))
	router.Path("/safari-pinned-tab.svg").Methods("GET").HandlerFunc(customHandler("favicon/safari-pinned-tab.svg", MimeTypeSVG))

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
