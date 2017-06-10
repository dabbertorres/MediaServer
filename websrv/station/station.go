package station

import (
	"bytes"
	"html/template"
	"net/http"
	"log"
	
	"MediaServer/internal/song"
	"github.com/gorilla/websocket"
	"MediaServer/websrv/msg"
	"MediaServer/urlgen"
)

type Playlist []song.Info

type Data struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Playlist `json:"playlist"`
	connections []Connection
	broadcast chan msg.Data
}

func New(w http.ResponseWriter, r *http.Request) Data {
	// TODO read station name from request json
	// Also maybe make .Name == .Url
	ret := Data{
		Name: "Station",
		Url: urlgen.Gen(),
		broadcast: make(chan msg.Data, 10),
	}
	
	ret.Add(w, r)
	
	// TODO launch a goroutine for handling station-wide messages
	
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
		Name: "placeholder", // JSON encoded in the request
		In: make(chan msg.Data, 10),
		Out: d.broadcast,
	})
	
	go d.connections[len(d.connections) - 1].Work(ws)
}

type Page struct {
	template *template.Template
	length int
}

func Load(data []byte) Page {
	return Page{
		template: template.Must(template.New("stationPage").Parse(string(data))),
		length: len(data),
	}
}

func (p Page) Generate(data Data) ([]byte, error) {
	ret := make([]byte, 0, p.length)
	err := p.template.Execute(bytes.NewBuffer(ret), data)
	return ret, err
}
