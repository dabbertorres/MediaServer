package msg

type StatusMessage struct {
	From         string `json:"from"`
	Disconnected bool   `json:"disconnected"`
}
