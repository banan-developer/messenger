package domain

type Message struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	FromID   int    `json:"from_id"`
	ToID     int    `json:"to_id"`
	CreateAt string `json:"created_at"`
	ChatId   int    `json:"chats_id"`
}

type ChatListItem struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	LastMessage     string `json:"last_message"`
	LastMessageTime string `json:"last_message_time"`
}
