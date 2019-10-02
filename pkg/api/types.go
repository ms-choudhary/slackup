package api

type Message struct {
	User     string    `json:"user"`
	Text     string    `json:"text"`
	Ts       string    `json:"ts"`
	Comments []Message `json:"comments"`
}
