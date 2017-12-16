package websrv

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	basePlaylistCapacity = 32
	broadcastBufferSize  = 8
)

type PlaybackUpdate int

const (
	SongOver PlaybackUpdate = iota
	SongStart
)

type ChatMessage struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

type StatusUpdate struct {
	Who  string `json:"who"`
	What string `json:"what"`
}

type Station struct {
	Name            string `json:"name"`
	Playlist        []Song `json:"playlist"`
	TuneIn          chan *Client
	TuneOut         chan *Client
	clients         map[*Client]bool
	playbackUpdates chan PlaybackUpdate
	chatUpdates     chan ChatMessage
	statusUpdates   chan StatusUpdate
}

func NewStation(name string) *Station {
	ret := &Station{
		Name:            name,
		Playlist:        make([]Song, 0, basePlaylistCapacity),
		TuneIn:          make(chan *Client, 2),
		TuneOut:         make(chan *Client, 2),
		clients:         make(map[*Client]bool),
		playbackUpdates: make(chan PlaybackUpdate, broadcastBufferSize),
		chatUpdates:     make(chan ChatMessage, broadcastBufferSize),
		statusUpdates:   make(chan StatusUpdate, broadcastBufferSize),
	}

	go ret.work()

	return ret
}

func (stn *Station) Add(w http.ResponseWriter, r *http.Request) {
	// TODO route handler should be doing this
	up := websocket.Upgrader{}
	ws, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// TODO get username from session sent with request
	c := NewClient("placeholder", ws)
	stn.clients[c] = true

	go c.read(stn)
	go c.write(stn)
}

func (stn *Station) work() {
	for {
		select {
		case client := <-stn.TuneIn:
			stn.clients[client] = true
			
		case client := <-stn.TuneOut:
			if _, ok := stn.clients[client]; ok {
				delete(stn.clients, client)
			}
			
		case play := <-stn.playbackUpdates:
			for c := range stn.clients {
				select {
				case c.playbackUpdates <- play:
				default:
					stn.TuneOut <- c
				}
			}
		
		case chat := <-stn.chatUpdates:
			for c := range stn.clients {
				select {
				case c.chatUpdates <- chat:
				default:
					stn.TuneOut <- c
				}
			}
		
		case sts := <-stn.statusUpdates:
			for c := range stn.clients {
				select {
				case c.statusUpdates <- sts:
				default:
					delete(stn.clients, c)
				}
			}
		}
	}
}
