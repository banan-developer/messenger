package repository

import (
	"database/sql"
	"log"
	"messenger_v2/internal/domain"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (m *MessageRepo) GetMessagesByChatID(FriendID, UserID int) ([]domain.Message, error) {
	var ChatID int
	err := m.db.QueryRow("SELECT id FROM chats WHERE (user1 = ? AND user2 = ?) OR (user1 = ? AND user2 = ?)", UserID, FriendID, FriendID, UserID).Scan(&ChatID)
	if err != nil {
		log.Println("БД: Ошибка при получении id чата")
		return nil, err
	}
	rows, err := m.db.Query("SELECT id, text, DATE_FORMAT(created_at, '%m-%d %H:%i') as created_at, from_id, to_id FROM messeges WHERE chats_id = ?", ChatID)

	if err != nil {
		log.Println("БД: Ошибка при получении сообщения")
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		rows.Scan(&message.Id, &message.Text, &message.CreateAt, &message.FromID, &message.ToID)
		message.ChatId = ChatID
		messages = append(messages, message)
	}
	if messages == nil {
		messages = []domain.Message{}
	}

	return messages, nil
}
