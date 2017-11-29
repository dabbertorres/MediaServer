package station

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	
	"radio/websrv"
	"radio/websrv/client"
)

const (
	basePlaylistCapacity = 32
	broadcastBufferSize  = 8
)

type PlaybackUpdate int

const (
	SongOver  PlaybackUpdate = iota
	SongStart
)

type ChatMessage struct {
	From string `json:"from"`
	Message string `json:"message"`
}

type StatusUpdate struct {
	Who string `json:"who"`
	What string `json:"what"`
}

type Map map[string]*Data

type Data struct {
	Name            string `json:"name"`
	Playlist        []websrv.Song `json:"playlist"`
	lock            sync.RWMutex
	clients         []client.Data
	playbackUpdates chan PlaybackUpdate
	chatUpdates     chan ChatMessage
	statusUpdates   chan StatusUpdate
}

func New(name string) *Data {
	ret := &Data{
		Name:            name,
		Playlist:        make([]websrv.Song, 0, basePlaylistCapacity),
		clients:         make([]client.Data, 0, 8),
		playbackUpdates: make(chan PlaybackUpdate, broadcastBufferSize),
		chatUpdates:     make(chan ChatMessage, broadcastBufferSize),
		statusUpdates:   make(chan StatusUpdate, broadcastBufferSize),
	}

	go ret.work()

	return ret
}

func (stn *Data) Add(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{}
	ws, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	stn.clients = append(stn.clients,
		client.Data{
			Name: "placeholder", // TODO json encoded in the request
		})

	go stn.clients[len(stn.clients)-1].Work(ws, stn)
}

func (stn *Data) work() {

}
