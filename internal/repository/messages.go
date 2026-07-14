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
	return &MessageRepo{db: db}
}

// GetOrCreateChat — найти или создать личный чат
func (r *MessageRepo) GetOrCreateChat(user1ID, user2ID int) (int, error) {
	var chatID int

	// Ищем существующий чат
	err := r.db.QueryRow(`
		SELECT c.id 
		FROM chats c
		JOIN users_has_chats uhc1 ON c.id = uhc1.chats_id AND uhc1.users_id = ?
		JOIN users_has_chats uhc2 ON c.id = uhc2.chats_id AND uhc2.users_id = ?
		WHERE c.is_group = 0
	`, user1ID, user2ID).Scan(&chatID)

	if err == nil {
		return chatID, nil
	}

	if err != sql.ErrNoRows {
		log.Println("GetOrCreateChat error:", err)
		return 0, err
	}

	// Создаём чат
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO chats (is_group, created_by) VALUES (0, ?)", user1ID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	chatID = int(id)

	// Добавляем участников
	_, err = tx.Exec(`
		INSERT INTO users_has_chats (chats_id, users_id) VALUES (?, ?), (?, ?)
	`, chatID, user1ID, chatID, user2ID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return chatID, nil
}

func (r *MessageRepo) GetChatsForUser(userID int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.is_group, COALESCE((SELECT other_member.users_id FROM users_has_chats other_member WHERE other_member.chats_id = c.id AND other_member.users_id <> ? LIMIT 1), 0),
			CASE WHEN c.is_group = 1 THEN COALESCE(c.title, '') ELSE COALESCE((
				SELECT u.name FROM users_has_chats other_member
				JOIN users u ON u.id = other_member.users_id
				WHERE other_member.chats_id = c.id AND other_member.users_id <> ? LIMIT 1
			), '') END AS name,
			CASE WHEN c.is_group = 1 THEN COALESCE(c.avatar_url, '') ELSE COALESCE((
				SELECT u.avatar_url FROM users_has_chats other_member
				JOIN users u ON u.id = other_member.users_id
				WHERE other_member.chats_id = c.id AND other_member.users_id <> ? LIMIT 1
			), '') END AS avatar,
			COALESCE((SELECT text FROM messeges WHERE chats_id = c.id ORDER BY created_at DESC LIMIT 1), '')
		FROM chats c
		JOIN users_has_chats uh ON uh.chats_id = c.id
		WHERE uh.users_id = ?
		ORDER BY c.id DESC
	`, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, otherUserID int
		var group bool
		var title, avatar, last string
		if err := rows.Scan(&id, &group, &otherUserID, &title, &avatar, &last); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{"id": id, "user_id": otherUserID, "is_group": group, "name": title, "avatar": avatar, "last_message": last})
	}
	return result, rows.Err()
}

// SaveMessage — сохранить сообщение
func (r *MessageRepo) SaveMessage(msg *domain.Message) error {
	result, err := r.db.Exec(`
		INSERT INTO messeges (text, from_id, to_id, chats_id, attachment_url) 
		VALUES (?, ?, ?, ?, ?)
	`, msg.Text, msg.FromID, msg.ToID, msg.ChatID, msg.AttachmentURL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	msg.ID = int(id)
	return nil
}

// GetMessagesByChatID — получить все сообщения чата
func (r *MessageRepo) GetMessagesByChatID(chatID int) ([]domain.Message, error) {
	if chatID <= 0 {
		return []domain.Message{}, nil
	}

	rows, err := r.db.Query(`
		SELECT id, text, DATE_FORMAT(created_at, '%m-%d %H:%i') as created_at, 
		       from_id, to_id, COALESCE(attachment_url, '') as attachment_url
		FROM messeges WHERE chats_id = ? ORDER BY created_at ASC
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]domain.Message, 0)
	for rows.Next() {
		var msg domain.Message
		err := rows.Scan(&msg.ID, &msg.Text, &msg.CreatedAt, &msg.FromID, &msg.ToID, &msg.AttachmentURL)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

// GetMessageByID — получить сообщение по ID
func (r *MessageRepo) GetMessageByID(msgID int) (*domain.Message, error) {
	var msg domain.Message
	err := r.db.QueryRow(`
		SELECT id, text, from_id, to_id, chats_id, attachment_url FROM messeges WHERE id = ?
	`, msgID).Scan(&msg.ID, &msg.Text, &msg.FromID, &msg.ToID, &msg.ChatID, &msg.AttachmentURL)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// UpdateMessage — обновить текст сообщения
func (r *MessageRepo) UpdateMessage(msgID int, newText string) error {
	if newText == "" {
		_, err := r.db.Exec("DELETE FROM messeges WHERE id = ?", msgID)
		return err
	}
	_, err := r.db.Exec("UPDATE messeges SET text = ? WHERE id = ?", newText, msgID)
	return err
}

// DeleteMessage — удалить сообщение
func (r *MessageRepo) DeleteMessage(msgID int) error {
	_, err := r.db.Exec("DELETE FROM messeges WHERE id = ?", msgID)
	return err
}

// GetChatParticipants — получить участников чата
func (r *MessageRepo) GetChatParticipants(chatID int) ([]int, error) {
	rows, err := r.db.Query(`
		SELECT users_id FROM users_has_chats WHERE chats_id = ?
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var userID int
		rows.Scan(&userID)
		users = append(users, userID)
	}
	return users, nil
}

func (r *MessageRepo) IsChatParticipant(chatID, userID int) (bool, error) {
	var exists int
	err := r.db.QueryRow("SELECT 1 FROM users_has_chats WHERE chats_id = ? AND users_id = ?", chatID, userID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
