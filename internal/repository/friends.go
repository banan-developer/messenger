package repository

import (
	"database/sql"
	"log"
	"messenger_v2/internal/domain"
)

type FriendRepo struct {
	db *sql.DB
}

func NewFriendRepo(db *sql.DB) *FriendRepo {
	return &FriendRepo{
		db: db,
	}
}

func (f *FriendRepo) GetFriendsByID(UserID int) ([]domain.FriendResponse, error) {
	rows, err := f.db.Query(
		`SELECT users.id, users.name, users.avatar_url
		FROM friends
		JOIN users ON users.id = friends.friend_id
		WHERE friends.users_id = ?
		AND friends.status = 'accepted';
		`, UserID)

	if err != nil {
		log.Println("БД: Ошибка при получении данных друзей")
		return nil, err
	}
	defer rows.Close()

	var friend []domain.FriendResponse
	for rows.Next() {
		var user domain.FriendResponse
		rows.Scan(&user.ID, &user.Name, &user.Avatar)
		friend = append(friend, user)
	}
	if friend == nil {
		friend = []domain.FriendResponse{}
	}

	return friend, nil
}

func (f *FriendRepo) AddToFriend(UserID int, FriendID int, Status string) error {
	_, err := f.db.Exec("INSERT INTO friends (friend_id, users_id, status) VALUES (?, ?, ?)", FriendID, UserID, Status)
	if err != nil {
		log.Println("БД: Ошибка занесение данных")
		return err
	}
	return nil
}

func (f *FriendRepo) FoundFriendByID(FriendName string) ([]domain.FriendResponse, error) {
	rows, err := f.db.Query("SELECT id, name, avatar_url FROM users WHERE name LIKE ?", "%"+FriendName+"%")

	if err != nil {
		log.Println("БД: Ошибка при получении данных друга")
		return nil, err
	}
	defer rows.Close()

	var friends []domain.FriendResponse
	for rows.Next() {
		var friend domain.FriendResponse
		rows.Scan(&friend.ID, &friend.Name, &friend.Avatar)
		friends = append(friends, friend)
	}
	if friends == nil {
		friends = []domain.FriendResponse{}
	}
	return friends, nil
}

func (f *FriendRepo) GetIncomingRequest(UserID int) ([]domain.FriendResponse, error) {
	rows, err := f.db.Query(
		`SELECT u.id, u.name, u.avatar_url
        FROM friends f
        JOIN users u ON u.id = f.users_id
        WHERE f.friend_id = ?
        AND f.status = 'invited'`,
		UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []domain.FriendResponse
	for rows.Next() {
		var friend domain.FriendResponse
		err := rows.Scan(&friend.ID, &friend.Name, &friend.Avatar)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	if friends == nil {
		friends = []domain.FriendResponse{}
	}

	return friends, nil
}

func (f *FriendRepo) AcceptComingRequset(FriendID, UserID int) error {
	tx, _ := f.db.Begin()
	defer tx.Rollback()

	tx.Exec(
		"UPDATE friends SET status = 'accepted' WHERE friend_id = ? AND users_id = ?",
		UserID, FriendID,
	)

	tx.Exec(
		"INSERT INTO friends (friend_id, users_id, status) VALUES (?, ?, 'accepted')",
		FriendID, UserID,
	)

	return tx.Commit()
}

func (f *FriendRepo) DeleteFriend(userID, friendID int) error {
	_, err := f.db.Exec("DELETE FROM friends WHERE (users_id=? AND friend_id=?) OR (users_id=? AND friend_id=?)", userID, friendID, friendID, userID)
	return err
}
