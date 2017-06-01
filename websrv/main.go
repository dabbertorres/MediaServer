package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	
	"MediaServer/websrv/file"
)

// TODO
type Station struct {

}

//var stations = make(map[string]Station)

func main() {
	if err := file.WatchInit(); err != nil {
		panic(err)
	}
	defer file.WatchStop()
	
	stop := make(chan bool)
	
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler("app/index.html", stop))
	serveDir(serverMux, "/html/", stop)
	serveDir(serverMux, "/css/", stop)
	serveDir(serverMux, "/js/", stop)
	
	// TODO request path for a media file
	// create request to the file db for the file.
	// it responds with (hopefully) the file
	// then we can stream that file!
	
	server := http.Server{
		Addr: ":8080",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 32,
		Handler: serverMux,
	}
	
	// listen for termination signal!
	go func() {
		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt)
		
		// block until we get something
		<-interrupt
		fmt.Println("Caught SIGINT")
		server.Close()
		close(stop)
	}()
	
	// hey now we can do what we want
	fmt.Println("Serving...")
	server.ListenAndServe()
	fmt.Println("Done")
}

func handler(filename string, stop <-chan bool) func(http.ResponseWriter, *http.Request) {
	event, err := file.Watch(filename, stop)
	if err != nil {
		panic(err)
	}
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	
	// we don't want to potentially block a response longer than needed,
	// so we gotta do file updates between requests.
	// which means we gotta mutex up the data
	mut := sync.Mutex{}
	
	// launch a goroutine which, on a good file watch event, updates ours data for this file
	go func() {
		for d := range event {
			if d.Error != nil {
				// on an error, log it, but don't terminate
				// the data we already had (should) still be good
				// so let's just stop watching the file
				fmt.Println(d.Error)
				return
			} else {
				mut.Lock()
				data = d.Data
				mut.Unlock()
			}
		}
	}()
	
	return func(w http.ResponseWriter, r *http.Request) {
		mut.Lock()
		w.Write(data)
		mut.Unlock()
	}
}

// pick a directory, and setup handlers for all files in there!
func serveDir(mux *http.ServeMux, dir string, stop <-chan bool) error {
	files, err := ioutil.ReadDir("app" + dir)
	if err != nil {
		return err
	}
	
	for _, f := range files {
		path := dir + f.Name()
		mux.HandleFunc(path, handler("app" + path, stop))
	}
	
	return nil
}
