package main

import "net/http"

type HandleFunc func(http.ResponseWriter, *http.Request)

var routes = map[string]HandleFunc{
	// basic server files!
	"/":      indexHandler,
	"/css/":  handler("text/css"),
	"/html/": handler("text/html"),
	"/js/":   handler("text/javascript"),

	// images!
	"/img/png": handler("image/png"),
	"/img/svg": handler("image/svg+xml"),

	// favicon config crap (bless you, realfavicongenerator.net)
	"/browserconfig.xml": customHandler(san("browserconfig.xml"), "application/xml"),
	"/manifest.json":     customHandler(san("manifest.json"), "application/json"),

	// now the actual favicons
	"/android-chrome-192x192.png": customHandler(san("img/favicon/android-chrome-192x192.png"), "image/png"),
	"/android-chrome-512x512.png": customHandler(san("img/favicon/android-chrome-512x512.png"), "image/png"),
	"/apple-touch-icon.png":       customHandler(san("img/favicon/apple-touch-icon.png"), "image/png"),
	"/favicon.ico":                customHandler(san("img/favicon/favicon.ico"), "image/x-icon"),
	"/favicon.png":                customHandler(san("img/favicon/favicon.png"), "image/png"),
	"/favicon-16x16.png":          customHandler(san("img/favicon/favicon-16x16.png"), "image/png"),
	"/favicon-32x32.png":          customHandler(san("img/favicon/favicon-32x32.png"), "image/png"),
	"/mstile-150x150.png":         customHandler(san("img/favicon/mstile-150x150.png"), "image/png"),
	"/safari-pinned-tab.svg":      customHandler(san("img/favicon/safari-pinned-tab.svg"), "image/svg"),

	// actually interesting stuff eventually
	"/song/":    songHandler,
	"/station/": stationHandler,
	"/search/":  searchHandler,
}

func makeMuxer() *http.ServeMux {
	mux := http.NewServeMux()

	for path, handle := range routes {
		mux.HandleFunc(path, handle)
	}

	return mux
}
