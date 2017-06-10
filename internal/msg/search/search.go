package search

import (
	"MediaServer/internal/song"
)

type Request struct {
	Title  string `json:"title,omitempty"`
	Album  string `json:"album,omitempty"`
	Artist string `json:"artist,omitempty"`
	Genre  string `json:"genre,omitempty"`
	Any    string `json:"any,omitempty"`
}

type Response []song.Info
