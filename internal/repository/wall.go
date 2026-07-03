package repository

import (
	"database/sql"
	"fmt"
	"log"
	"messenger_v2/internal/domain"
)

type WallRepo struct {
	db *sql.DB
}

func NewWallRepo(db *sql.DB) *WallRepo {
	return &WallRepo{
		db: db,
	}
}

func (w *WallRepo) CreatePost(res *domain.CreateWallRequest, userID int) error {
	_, err := w.db.Exec("INSERT INTO wall (title, text, users_id, img_scr) VALUES (?, ?, ?, ?)", res.Title, res.Text, userID, res.Img)
	if err != nil {
		log.Printf("Ошибка получении поста %w", err)
		return err
	}
	return nil
}

func (w *WallRepo) GetPost(UserID int) ([]domain.WallPost, error) {
	rows, err := w.db.Query("SELECT idwall, title, text, img_scr FROM wall WHERE users_id  = ?", UserID)
	if err != nil {
		log.Printf("Ошибка получении поста %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []domain.WallPost

	for rows.Next() {
		var post domain.WallPost
		rows.Scan(&post.Id, &post.Title, &post.Text, &post.Img)
		posts = append(posts, post)
		fmt.Println(post.Img)
	}

	if posts == nil {
		posts = []domain.WallPost{}
	}

	return posts, nil
}
