package repository

import (
	"database/sql"
	"fmt"
	"messenger_v2/internal/domain"
)

type ChatRepo struct {
	db *sql.DB
}

func (r *ChatRepo) GetGroupsForUser(userID int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`SELECT c.id, c.title, COALESCE(c.avatar_url,'') FROM chats c JOIN users_has_chats u ON u.chats_id=c.id WHERE u.users_id=? AND c.is_group=1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	groups := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var title, avatar string
		if err := rows.Scan(&id, &title, &avatar); err != nil {
			return nil, err
		}
		groups = append(groups, map[string]interface{}{"id": id, "title": title, "avatar_url": avatar})
	}
	return groups, rows.Err()
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) IsParticipant(chatID, userID int) (bool, error) {
	var value int
	err := r.db.QueryRow("SELECT 1 FROM users_has_chats WHERE chats_id = ? AND users_id = ?", chatID, userID).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *ChatRepo) IsCreator(chatID, userID int) (bool, error) {
	var value int
	err := r.db.QueryRow("SELECT 1 FROM chats WHERE id = ? AND created_by = ? AND is_group = 1", chatID, userID).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *ChatRepo) AreFriends(userID, friendID int) (bool, error) {
	var value int
	err := r.db.QueryRow(`SELECT 1 FROM friends WHERE users_id = ? AND friend_id = ? AND status = 'accepted' LIMIT 1`, userID, friendID).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// CreateGroupChat — создать группу
func (r *ChatRepo) CreateGroupChat(title, avatarURL string, userIDs []int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO chats (is_group, title, avatar_url, created_by) VALUES (1, ?, ?, ?)
	`, title, avatarURL, userIDs[0])
	if err != nil {
		return 0, err
	}
	chatID64, _ := result.LastInsertId()
	chatID := int(chatID64)

	if len(userIDs) == 0 || len(userIDs) > 50 {
		return 0, fmt.Errorf("group must contain from 1 to 50 members")
	}
	for _, userID := range userIDs {
		_, err = tx.Exec(`
			INSERT INTO users_has_chats (chats_id, users_id) VALUES (?, ?)
		`, chatID, userID)
		if err != nil {
			return 0, err
		}
	}

	tx.Commit()
	return chatID, nil
}

// GetGroupMembers — получить участников группы
func (r *ChatRepo) GetGroupMembers(chatID int) ([]domain.User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.name, u.avatar_url
		FROM users_has_chats uhc
		JOIN users u ON uhc.users_id = u.id
		WHERE uhc.chats_id = ?
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		rows.Scan(&user.ID, &user.Name, &user.Avatar)
		users = append(users, user)
	}
	return users, nil
}

// AddMemberToGroup — добавить участника
func (r *ChatRepo) AddMemberToGroup(chatID, userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO users_has_chats (chats_id, users_id)
		SELECT ?, ? WHERE (SELECT COUNT(*) FROM users_has_chats WHERE chats_id = ?) < 50
	`, chatID, userID, chatID)
	return err
}

// RemoveMemberFromGroup — удалить участника
func (r *ChatRepo) RemoveMemberFromGroup(chatID, userID int) error {
	var ownerID, memberCount int
	if err := r.db.QueryRow(`SELECT created_by, (SELECT COUNT(*) FROM users_has_chats WHERE chats_id = ?) FROM chats WHERE id = ? AND is_group = 1`, chatID, chatID).Scan(&ownerID, &memberCount); err != nil {
		return err
	}
	if userID == ownerID {
		return fmt.Errorf("group creator cannot be removed")
	}
	if memberCount <= 1 {
		return fmt.Errorf("last group member cannot be removed")
	}
	_, err := r.db.Exec(`
		DELETE FROM users_has_chats WHERE chats_id = ? AND users_id = ?
	`, chatID, userID)
	return err
}

func (r *ChatRepo) RenameGroup(chatID int, title string) error {
	_, err := r.db.Exec("UPDATE chats SET title = ? WHERE id = ? AND is_group = 1", title, chatID)
	return err
}
func (r *ChatRepo) DeleteGroup(chatID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("DELETE FROM messeges WHERE chats_id = ?", chatID); err != nil {
		return err
	}
	if _, err = tx.Exec("DELETE FROM users_has_chats WHERE chats_id = ?", chatID); err != nil {
		return err
	}
	if _, err = tx.Exec("DELETE FROM chats WHERE id = ? AND is_group = 1", chatID); err != nil {
		return err
	}

	return tx.Commit()
}
