package api

type Message struct {
	User     string    `json:"user"`
	Text     string    `json:"text"`
	Ts       string    `json:"ts"`
	Comments []Message `json:"comments"`
}

func (m *Message) UpdateMessage(user, text, ts string) {
	m.User = user
	m.Text = text
	m.Ts = ts
}

func (m *Message) AddComment(c Message) {
	m.Comments = append(m.Comments, c)
}
