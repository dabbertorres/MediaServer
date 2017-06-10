// messages for communicating between a station on the server and clients
package msg

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type Type int

const (
	TypeChat Type = iota
	TypeStream
	TypePlaylist
	TypeStatus
)

// implements json.Unmarshaler
type Data struct {
	Type Type `json:"type"`

	Chat     *ChatMessage     `json:"chat,omitempty"`
	Stream   *StreamMessage   `json:"stream,omitempty"`
	Playlist *PlaylistMessage `json:"playlist,omitempty"`
	Status   *StatusMessage   `json:"status,omitempty"`
}

func (m Data) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return json.UnmarshalTypeError{
			Value:  "null",
			Type:   reflect.TypeOf((*Type)(nil)),
			Offset: 0,
			Struct: "websrv.msg.Data",
			Field:  "Type",
		}
	}

	dec := json.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&m.Type); err != nil {
		return err
	}

	if m.Type < TypeChat || m.Type > TypeStatus {
		return json.UnmarshalTypeError{
			Value:  string(m.Type),
			Type:   reflect.TypeOf((*Type)(nil)),
			Struct: "websrv.msg.Data",
			Field:  "Type",
		}
	}

	var iv interface{}
	switch m.Type {
	case TypeChat:
		m.Chat = &ChatMessage{}
		iv = m.Chat

	case TypeStream:
		m.Stream = &StreamMessage{}
		iv = m.Stream

	case TypePlaylist:
		m.Playlist = &PlaylistMessage{}
		iv = m.Playlist

	case TypeStatus:
		m.Status = &StatusMessage{}
		iv = m.Status
	}

	if err := dec.Decode(iv); err != nil {
		return err
	}

	return nil
}
