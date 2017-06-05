package media

import (
	"io"
	"io/ioutil"
	"time"
)

type SongInfo struct {
	Title  string        `json:"title"`
	Artist string        `json:"artist"`
	Album  string        `json:"album"`
	Length time.Duration `json:"length"`
}

type SongData []byte

type Song struct {
	SongInfo
	data SongData
}

// NextChunk returns a data slice, and the amount of bytes left total,
// given a starting point, and a max amount to retrieve
// if start is out of bounds, (nil, -1) is returned.
func (s *Song) NextChunk(start int, size int) ([]byte, int) {
	dataLen := len(s.data)
	if start >= dataLen {
		return nil, -1
	}

	end := start + size
	left := dataLen - start
	if end > left {
		end = left
	}

	return s.data[start:end], dataLen - end
}

func (s *Song) Load(buf io.Reader) (err error) {
	s.data, err = ioutil.ReadAll(buf)
	return
}
