package service

import (
	"errors"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepo
}

func NewMessageService(repo *repository.MessageRepo) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) GetOrCreateChat(user1ID, user2ID int) (int, error) {
	if user1ID <= 0 || user2ID <= 0 {
		return 0, errors.New("invalid user ids")
	}
	return s.repo.GetOrCreateChat(user1ID, user2ID)
}

func (s *MessageService) SaveMessage(msg *domain.Message) error {
	// ✅ теперь проверяем текст ИЛИ вложение
	if msg.Text == "" && msg.AttachmentURL == "" {
		return errors.New("message cannot be empty")
	}
	if msg.FromID <= 0 || msg.ChatID <= 0 {
		return errors.New("invalid sender or chat")
	}
	return s.repo.SaveMessage(msg)
}

func (s *MessageService) GetMessagesByChatID(chatID int) ([]domain.Message, error) {
	if chatID <= 0 {
		return nil, errors.New("invalid chat id")
	}
	return s.repo.GetMessagesByChatID(chatID)
}

func (s *MessageService) GetMessageByID(msgID int) (*domain.Message, error) {
	if msgID <= 0 {
		return nil, errors.New("invalid message id")
	}
	return s.repo.GetMessageByID(msgID)
}

func (s *MessageService) UpdateMessage(msgID int, newText string) error {
	if msgID <= 0 {
		return errors.New("invalid message id")
	}
	return s.repo.UpdateMessage(msgID, newText)
}

func (s *MessageService) DeleteMessage(msgID int) error {
	if msgID <= 0 {
		return errors.New("invalid message id")
	}
	return s.repo.DeleteMessage(msgID)
}

func (s *MessageService) GetChatsForUser(userID int) ([]map[string]interface{}, error) { return s.repo.GetChatsForUser(userID) }

func (s *MessageService) IsChatParticipant(chatID, userID int) (bool, error) {
	return s.repo.IsChatParticipant(chatID, userID)
}
