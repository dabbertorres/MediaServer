package station

import (
	"log"
	"net/http"

	"MediaServer/internal/song"
	"MediaServer/websrv/msg"
	"github.com/gorilla/websocket"
)

type Data struct {
	Name        string      `json:"name"`
	Playlist    []song.Info `json:"playlist"`
	connections []Connection
	broadcast   chan msg.Data
}

func New(w http.ResponseWriter, r *http.Request) Data {
	// TODO read station name from request json
	// Also maybe make .Name == .Url
	ret := Data{
		Name:      "Station",
		broadcast: make(chan msg.Data, 10),
	}

	// TODO launch a goroutine for handling station-wide messages

	return ret
}

func (d *Data) Add(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{}
	ws, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	d.connections = append(d.connections, Connection{
		Name: "placeholder", // TODO json encoded in the request
		In:   make(chan msg.Data, 10),
		Out:  d.broadcast,
	})

	go d.connections[len(d.connections)-1].Work(ws)
}
