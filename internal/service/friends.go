package service

import (
	"errors"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
)

type FriendService struct {
	repo *repository.FriendRepo
}

func NewFrinedService(repo *repository.FriendRepo) *FriendService {
	return &FriendService{
		repo: repo,
	}
}

func (f *FriendService) GetFriendsByID(UserID int) ([]domain.FriendResponse, error) {
	if UserID <= 0 {
		log.Println("Пользователь не найден")
		return nil, errors.New("Пользователь не найден")
	}
	return f.repo.GetFriendsByID(UserID)
}

func (f *FriendService) AddToFriend(UserID int, FriendID int, Status string) error {
	if UserID <= 0 {
		log.Println("Пользователь не найден")
		return errors.New("Пользователь не найден")
	}
	if FriendID <= 0 {
		log.Println("Пользователь не найден")
		return errors.New("Пользователь не найден")
	}

	if Status == "" {
		log.Println("Операция не найдена")
		return errors.New("Операция не найдена")
	}
	return f.repo.AddToFriend(UserID, FriendID, Status)

}

func (f *FriendService) FoundFriendByID(FriendName string) ([]domain.FriendResponse, error) {
	if FriendName == "" {
		log.Println("Пользователь не найден")
		return nil, errors.New("Пользователь не найден")
	}
	return f.repo.FoundFriendByID(FriendName)
}

func (f *FriendService) GetIncomingRequest(UserID int) ([]domain.FriendResponse, error) {
	if UserID <= 0 {
		log.Println("Пользователь не найден")
		return nil, errors.New("Пользователь не найден")
	}
	return f.repo.GetIncomingRequest(UserID)
}

func (f *FriendService) AcceptComingRequset(FriendID, UserID int) error {
	if UserID <= 0 {
		return errors.New("Пользователь не найден")
	}
	if FriendID <= 0 {
		return errors.New("Пользователь(Друг) не найден")
	}
	return f.repo.AcceptComingRequset(FriendID, UserID)
}
