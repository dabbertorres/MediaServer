package websrv

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

var (
	Disconnecting = errors.New("client disconnecting")
)

type MessageType int

const (
	MessageTypePlayback MessageType = iota
	MessageTypeChat
	MessageTypeStatus
)

type WebSocketMessage struct {
	Type MessageType `json:"type"`
	Msg  map[string]interface{} `json:"msg"`
}

type Client struct {
	Name            string
	playbackUpdates chan PlaybackUpdate
	chatUpdates     chan ChatMessage
	statusUpdates   chan StatusUpdate
	ws              *websocket.Conn
}

func NewClient(name string, ws *websocket.Conn) *Client {
	return &Client{
		Name:            name,
		playbackUpdates: make(chan PlaybackUpdate),
		chatUpdates:     make(chan ChatMessage),
		statusUpdates:   make(chan StatusUpdate),
		ws:              ws,
	}
}

func (c *Client) Close() {
	close(c.playbackUpdates)
	close(c.chatUpdates)
	close(c.statusUpdates)
	c.ws.Close()
}

func (c *Client) read(stn *Station) {
	defer func() {
		stn.TuneOut <- c
		c.Close()
	}()

	for {
		msg := WebSocketMessage{}
		if err := c.ws.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Unexpected close error in client read loop:", err)
			}
			return
		}

		switch msg.Type {
		case MessageTypePlayback:
			play, ok := msg.Msg["update"].(PlaybackUpdate)
			if !ok {
			
			}
			
			stn.playbackUpdates <- play

		case MessageTypeChat:

		case MessageTypeStatus:
		}
	}
}

func (c *Client) write(stn *Station) {

}

// TODO do some filtering of error types. Not all errors are fatal for the connection.
func (c *Client) handle(fromClient bool, m Client, ws *websocket.Conn) error {
	/*switch m.Type {
	case TypeChat:
		if fromClient {
			c.Out <- m
		} else {
			if err := ws.WriteJSON(&m); err != nil {
				log.Println("ChatMessage (Server) error:", err)
				return err
			}
		}

	case TypeStream:
		if fromClient {
			c.Out <- m
		} else {
		}

	case TypePlaylist:
		if fromClient {
			c.Out <- m
		} else {
			err := ws.WriteJSON(&m)
			if err != nil {
				log.Println("PlaylistMessage (Server) error:", err)
			}
		}

	case TypeStatus:
		if fromClient {
			c.Out <- m
			if m.Status.Disconnected {
				return Disconnecting
			}
		} else {
			// let client know who disconnected
			err := ws.WriteJSON(&Client{
				Type: TypeChat,
				Chat: &ChatMessage{
					From:    "Station",
					Content: m.Status.From + " disconnected.",
				},
			})
			if err != nil {
				log.Println("StatusUpdate (Server) error:", err)
				return err
			}
		}
	}*/

	return nil
}
