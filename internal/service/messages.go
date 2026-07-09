package service

import (
	"errors"
	"fmt"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepo
}

func NewMessagesService(repo *repository.MessageRepo) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (m *MessageService) GetMessagesByChatID(FriendID, UserID int) ([]domain.Message, error) {
	if FriendID <= 0 {
		log.Println("Друг не найден")
		return nil, errors.New("Друг не найден")
	}
	if UserID <= 0 {
		fmt.Println(UserID)
		log.Println("Пользователь не найден")
		return nil, errors.New("Пользователь не найден")
	}
	return m.repo.GetMessagesByChatID(FriendID, UserID)
}
