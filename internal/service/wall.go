package service

import (
	"errors"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
)

type WallService struct {
	repo *repository.WallRepo
}

func NewWallService(repo *repository.WallRepo) *WallService {
	return &WallService{
		repo: repo,
	}
}

func (r *WallService) CreatePost(res *domain.CreateWallRequest, UserID int) ([]domain.WallPost, error) {
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

	return r.repo.GetPost(UserID)
}

func (r *WallService) GetPost(UserId int) ([]domain.WallPost, error) {
	if UserId <= 0 {
		return nil, errors.New("Пользователь не найден")
	}
	return r.repo.GetPost(UserId)
}
