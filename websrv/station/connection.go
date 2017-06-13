package station

import (
	"github.com/gorilla/websocket"
	"MediaServer/websrv/msg"
	"log"
	"errors"
)

var (
	Disconnecting = errors.New("Client disconnecting")
)

type Connection struct {
	Name string
	In chan msg.Data
	Out chan<- msg.Data
}

func (conn *Connection) Work(ws *websocket.Conn) {
	defer ws.Close()
	
	// workaround for websocket.Conn not having a message reading channel!
	type WsMsg struct {
		msg.Data
		Error error
	}
	wsMsg := make(chan WsMsg)
	
	go func() {
		for {
			m := msg.Data{}
			err := ws.ReadJSON(&m)
			wsMsg <- WsMsg{m, err}
		}
	}()
	
	for {
		select {
		// messages from the client
		case raw := <-wsMsg:
			if raw.Error != nil && !websocket.IsUnexpectedCloseError(raw.Error, websocket.CloseGoingAway) {
				log.Println("WebSocket error:", raw.Error)
				return
			}
			
			err := conn.handle(true, raw.Data, ws)
			if err == Disconnecting {
				log.Println(err)
				return
			} else if err != nil {
				log.Println("Message error:", err)
				return
			}
			
		// messages from the server/station
		case m := <-conn.In:
			err := conn.handle(false, m, ws)
			if err != nil {
				log.Println("Message error:", err)
				return
			}
		}
	}
}

// TODO do some filtering of error types. Not all errors are fatal for the connection.
func (conn *Connection) handle(fromClient bool, m msg.Data, ws *websocket.Conn) error {
	switch m.Type {
	case msg.TypeChat:
		if fromClient {
			conn.Out <- m
		} else {
			if err := ws.WriteJSON(&m); err != nil {
				log.Println("ChatMessage (Server) error:", err)
				return err
			}
		}
	
	case msg.TypeStream:
		if fromClient {
			conn.Out <- m
		} else {
		}
	
	case msg.TypePlaylist:
		if fromClient {
			conn.Out <- m
		} else {
			err := ws.WriteJSON(&m)
			if err != nil {
				log.Println("PlaylistMessage (Server) error:", err)
			}
		}
	
	case msg.TypeStatus:
		if fromClient {
			conn.Out <- m
			if m.Status.Disconnected {
				return Disconnecting
			}
		} else {
			// let client know who disconnected
			err := ws.WriteJSON(&msg.Data{
				Type: msg.TypeChat,
				Chat: &msg.ChatMessage{
					From: "Station",
					Content: m.Status.From + " disconnected.",
				},
			})
			if err != nil {
				log.Println("StatusMessage (Server) error:", err)
				return err
			}
		}
	}
	
	return nil
}
