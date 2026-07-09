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

func (m *MessageRepo) GetChatsWithLastMessages(UserID int) ([]domain.ChatListItem, error) {
	rows, err := m.db.Query(`
        SELECT 
            c.id,
            u.id AS user_id,
            u.name,
            u.avatar_url,
            COALESCE(m.text, '') AS last_message,
            COALESCE(m.created_at, '') AS last_message_time
        FROM chats c
        JOIN users u ON u.id = CASE 
            WHEN c.user1 = ? THEN c.user2
            WHEN c.user2 = ? THEN c.user1
        END
        LEFT JOIN messeges m ON m.chats_id = c.id
        AND m.created_at = (
            SELECT MAX(created_at) FROM messeges WHERE chats_id = c.id
        )
        WHERE c.user1 = ? OR c.user2 = ?
        ORDER BY m.created_at DESC
    `, UserID, UserID, UserID, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatListItem
	for rows.Next() {
		var message domain.ChatListItem
		rows.Scan(&message.ID, &message.UserID, &message.Name, &message.Avatar, &message.LastMessage, &message.LastMessageTime)
		messages = append(messages, message)
	}
	if messages == nil {
		messages = []domain.ChatListItem{}
	}

	return messages, nil
}
