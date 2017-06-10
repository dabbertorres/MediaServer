package song

import (
	"io"
	"io/ioutil"
)

type Info struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Album  string  `json:"album"`
	Length float32 `json:"length,omitempty"` // time in seconds
}

type Data []byte

// NextChunk returns a data slice, and the amount of bytes left total,
// given a starting point, and a max amount to retrieve
// if start is out of bounds, (nil, -1) is returned.
func (data Data) NextChunk(start int, size int) ([]byte, int) {
	dataLen := len(data)
	if start >= dataLen {
		return nil, -1
	}

	end := start + size
	left := dataLen - start
	if end > left {
		end = left
	}

	return data[start:end], dataLen - end
}

func Load(buf io.Reader) (data Data, err error) {
	data, err = ioutil.ReadAll(buf)
	return
}
