package msg

type ChatMessage struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
