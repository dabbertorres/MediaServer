package msg

import (
	"bytes"
	"encoding/json"
	"reflect"
    
	"radio/internal/song"
)

type PlaylistMethod int

const (
	// from client
	MethodAppend PlaylistMethod = iota
	MethodPrepend
	MethodRemove
	MethodReorder

	// from server
	MethodUpdate
)

type PlaylistMessage struct {
	Method PlaylistMethod `json:"method"`

	Append  *PlaylistAppend  `json:"append,omitempty"`
	Prepend *PlaylistPrepend `json:"prepend,omitempty"`
	Remove  *PlaylistRemove  `json:"remove,omitempty"`
	Reorder *PlaylistReorder `json:"reorder,omitempty"`
	Update  *PlaylistUpdate  `json:"update,omitempty"`
}

func (pm PlaylistMessage) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return &json.UnmarshalTypeError{
			Value:  "null",
			Type:   reflect.TypeOf((*PlaylistMethod)(nil)),
			Offset: 0,
			Struct: "websrv.msg.PlaylistMessage",
			Field:  "Method",
		}
	}

	dec := json.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&pm.Method); err != nil {
		return err
	}

	if pm.Method < MethodAppend || pm.Method > MethodUpdate {
		return &json.UnmarshalTypeError{
			Value:  string(pm.Method),
			Type:   reflect.TypeOf((*PlaylistMethod)(nil)),
			Struct: "websrv.msg.PlaylistMessage",
			Field:  "Method",
		}
	}

	var iv interface{}
	switch pm.Method {
	case MethodAppend:
		pm.Append = &PlaylistAppend{}
		iv = pm.Append

	case MethodPrepend:
		pm.Prepend = &PlaylistPrepend{}
		iv = pm.Prepend

	case MethodRemove:
		pm.Remove = &PlaylistRemove{}
		iv = pm.Remove

	case MethodReorder:
		pm.Reorder = &PlaylistReorder{}
		iv = pm.Reorder

	case MethodUpdate:
		pm.Update = &PlaylistUpdate{}
		iv = pm.Update
	}

	if err := dec.Decode(iv); err != nil {
		return err
	}

	return nil
}

type PlaylistAppend []song.Info
type PlaylistPrepend []song.Info
type PlaylistRemove []int
type PlaylistReorder []struct {
	From int `json:"from"`
	To   int `json:"to"`
}
type PlaylistUpdate []song.Info
