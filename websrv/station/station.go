package station

import (
	"log"
	"net/http"
	"sync"

	"MediaServer/internal/song"
	"MediaServer/websrv/msg"
	"github.com/gorilla/websocket"
)

type Data struct {
	Name        string      `json:"name"`
	Playlist    []song.Info `json:"playlist"`
	lock        sync.RWMutex
	connections []Connection
	broadcast   chan msg.Data
}

func New(name string) *Data {
	const (
		basePlaylistCapacity = 32
		broadcastBufferSize = 16
	)
	
	ret := &Data{
		Name:      name,
		Playlist: make([]song.Info, 0, basePlaylistCapacity),
		broadcast: make(chan msg.Data, broadcastBufferSize),
	}
	
	// TODO TEST
	for i := 0; i < 40; i++ {
		ret.Playlist = append(ret.Playlist, song.Info{
			Artist: "Dead Sara",
			Album:  "Pleasure to Meet You",
			Title:  "Radio One Two",
			Length: 3*60 + 41,
		})
	}
	// TODO TEST

	go work(ret)

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

	go d.connections[len(d.connections)-1].Work(ws, d)
}

func work(data *Data) {

}
