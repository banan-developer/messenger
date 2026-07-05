package service

import (
	"errors"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
)

type PostService struct {
	repo *repository.WallRepo
}

func NewWallService(repo *repository.WallRepo) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (r *PostService) CreatePost(res *domain.CreateWallRequest, UserID int) ([]domain.WallPost, error) {
	if len(res.Title) <= 0 {
		return nil, errors.New("Название поста не может быть пустым")
	}
	if len(res.Text) <= 0 {
		return nil, errors.New("Пост не может быть пустым")
	}

	err := r.repo.CreatePost(res, UserID)
	if err != nil {
		return nil, errors.New("Ошибка при создании поста")
	}

	return r.repo.GetPostsByUserID(UserID)
}

func (r *PostService) GetPost(UserId int) ([]domain.WallPost, error) {
	if UserId <= 0 {
		return nil, errors.New("Пользователь не найден")
	}
	return r.repo.GetPostsByUserID(UserId)
}

func (r *PostService) DeletePost(PostID, UserID int) error {
	if PostID <= 0 {
		log.Println("Ошибка: ID поста не может быть отрицательным")
		return errors.New("Пост не может быть отрицательным")
	}
	if UserID <= 0 {
		return errors.New("Пользователь не найден")
	}
	return r.repo.DeletePost(PostID, UserID)
}

func (r *PostService) EditPostByID(res *domain.CreateWallRequest, PostID int) error {
	if PostID <= 0 {
		return errors.New("Пост не может быть отрицательным")
	}
	if res.Text == "" {
		return errors.New("Текст пост не может быть пустым")
	}
	if res.Title == "" {
		return errors.New("Заголовок поста не может быть пустым")
	}
	return r.repo.EditPostById(res, PostID)

}
