package media

type PlayList []SongInfo

type Station struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	PlayList `json:"playlist"`
}
