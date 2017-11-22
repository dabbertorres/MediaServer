package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"radio/internal/msg/search"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)

	searchParams := search.Request{}
	if err := json.Unmarshal(data, &searchParams); err != nil {
		log.Printf("Bad search request: '%s'\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func streamHandler(w http.ResponseWriter, r *http.Request) {

}
