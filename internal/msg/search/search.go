package search

import (
	"MediaServer/internal/media"
)

type Type string

const (
	Title  Type = "title"
	Album  Type = "album"
	Artist Type = "artist"
	Genre  Type = "genre"
	Any    Type = "any"
)

type Request map[Type]string

type Result []media.SongInfo
