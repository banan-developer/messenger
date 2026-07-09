package service

import (
	"errors"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(userID int) (*domain.User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetUserById(userID)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	if user.ID <= 0 {
		return errors.New("invalid user id")
	}
	if user.Name == "" {
		return errors.New("invalid user name")
	}
	if len(user.Name) > 100 {
		return errors.New("ivalid len of name")
	}

	return s.repo.UpdateUser(user)
}

func (s *UserService) UploadAvatarUser(UserID int, AvatarURL string) error {
	if UserID <= 0 {
		return errors.New("invalid user id")
	}
	if AvatarURL == "" {
		return errors.New("invalid avatarURL")
	}
	return s.repo.UploadAvatarUser(UserID, AvatarURL)
}

func (s *UserService) GetPersonByID(PersonID int) (*domain.UserResponse, error) {
	if PersonID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetPersonByID(PersonID)
}
