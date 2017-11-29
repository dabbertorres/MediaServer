package client

import (
	"errors"
	"log"
 
	"github.com/gorilla/websocket"
 
	"radio/websrv/station"
)

var (
	Disconnecting = errors.New("client disconnecting")
)

type Data struct {
	Name string
}

func (conn *Data) Work(ws *websocket.Conn, stn *station.Data) {
	defer ws.Close()

	// workaround for websocket.Conn not having a message reading channel!

	go func() {
		for {
			m := Data{}
			err := ws.ReadJSON(&m)
			wsMsg <- WsMsg{m, err}
		}
	}()
	
	// provide client with the current playlist
	stn.lock.RLock()
	ws.WriteJSON(&Data{
		Type: TypePlaylist,
		Playlist: &PlaylistMessage{
			Method: MethodUpdate,
			Update: (*PlaylistUpdate)(&data.Playlist),
		},
	})
	stn.lock.RUnlock()

	for {
		select {
		// messages from the client
		case raw := <-wsMsg:
			if raw.Error != nil && websocket.IsUnexpectedCloseError(raw.Error, websocket.CloseGoingAway) {
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
func (conn *Data) handle(fromClient bool, m Data, ws *websocket.Conn) error {
	switch m.Type {
	case TypeChat:
		if fromClient {
			conn.Out <- m
		} else {
			if err := ws.WriteJSON(&m); err != nil {
				log.Println("ChatMessage (Server) error:", err)
				return err
			}
		}

	case TypeStream:
		if fromClient {
			conn.Out <- m
		} else {
		}

	case TypePlaylist:
		if fromClient {
			conn.Out <- m
		} else {
			err := ws.WriteJSON(&m)
			if err != nil {
				log.Println("PlaylistMessage (Server) error:", err)
			}
		}

	case TypeStatus:
		if fromClient {
			conn.Out <- m
			if m.Status.Disconnected {
				return Disconnecting
			}
		} else {
			// let client know who disconnected
			err := ws.WriteJSON(&Data{
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
	}

	return nil
}
