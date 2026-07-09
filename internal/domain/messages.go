package domain

type Message struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	FromID   int    `json:"from_id"`
	ToID     int    `json:"to_id"`
	CreateAt string `json:"created_at"`
	ChatId   int    `json:"chats_id"`
}
