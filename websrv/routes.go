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

	// application file paths

	router.Path("/").Methods("GET").HandlerFunc(templateHandler("index.html", MimeTypeHTML))
	router.Path("/{path}.css").Methods("GET").HandlerFunc(handler(MimeTypeCSS))
	router.Path("/{path}.js").Methods("GET").HandlerFunc(handler(MimeTypeJavaScript))

	// favicons

	router.Path("/browserconfig.xml").Methods("GET").HandlerFunc(pathHandler("favicon/browserconfig.xml", MimeTypeXML))
	router.Path("/manifest.json").Methods("GET").HandlerFunc(pathHandler("favicon/manifest.json", MimeTypeJSON))
	router.Path("/android-chrome-192x192.png").Methods("GET").HandlerFunc(pathHandler("favicon/android-chrome-192x192.png", MimeTypePNG))
	router.Path("/android-chrome-512x512.png").Methods("GET").HandlerFunc(pathHandler("favicon/android-chrome-512x512.png", MimeTypePNG))
	router.Path("/apple-touch-icon.png").Methods("GET").HandlerFunc(pathHandler("favicon/apple-touch-icon.png", MimeTypePNG))
	router.Path("/favicon.ico").Methods("GET").HandlerFunc(pathHandler("favicon/favicon.ico", MimeTypeICO))
	router.Path("/favicon.png").Methods("GET").HandlerFunc(pathHandler("favicon/favicon.png", MimeTypePNG))
	router.Path("/favicon-16x16.png").Methods("GET").HandlerFunc(pathHandler("favicon/favicon-16x16.png", MimeTypePNG))
	router.Path("/favicon-32x32.png").Methods("GET").HandlerFunc(pathHandler("favicon/favicon-32x32.png", MimeTypePNG))
	router.Path("/mstile-150x150.png").Methods("GET").HandlerFunc(pathHandler("favicon/mstile-150x150.png", MimeTypePNG))
	router.Path("/safari-pinned-tab.svg").Methods("GET").HandlerFunc(pathHandler("favicon/safari-pinned-tab.svg", MimeTypeSVG))

	// api paths

	router.Path("/search").Methods("GET").HandlerFunc(searchHandler)

	router.Path("/song/{artist}/{title}").Methods("GET").HandlerFunc(songHandler)

	station := router.PathPrefix("/station").Subrouter()
	station.Path("/{station}").Methods("GET").HandlerFunc(tuneToStation)
	station.Path("/{station}").Methods("POST").HandlerFunc(startStation)
	station.Path("/{station}/playlist").Methods("GET").HandlerFunc(stationGetPlaylist)
	station.Path("/{station}/playlist").Methods("POST").HandlerFunc(stationAddSong)
	station.Path("/{station}/playlist").Methods("DELETE").HandlerFunc(stationRemoveSong)
	station.Path("/{station}/song").Methods("GET").HandlerFunc(stationGetPlayingSong)
	station.Path("/{station}/socket").Methods("GET").HandlerFunc(stationConnect)

	return router
}
