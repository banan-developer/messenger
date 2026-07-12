package domain

type Message struct {
	ID            int    `json:"id"`
	Text          string `json:"text"`
	FromID        int    `json:"from_id"`
	ToID          int    `json:"to_id"`
	ChatID        int    `json:"chat_id"`
	CreatedAt     string `json:"created_at"`
	AttachmentURL string `json:"attachment_url,omitempty"`
}

type EditMessageRequest struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
