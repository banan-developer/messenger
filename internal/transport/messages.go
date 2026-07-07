package transport

import (
	"encoding/json"
	"messenger_v2/internal/service"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

func (m *MessageHandler) Messages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.GetMessagesByChatID(w, r)
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func (m *MessageHandler) GetMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("id")
	friendID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	userIdParam := r.URL.Query().Get("user_id")
	UserID, _ := strconv.Atoi(userIdParam)

	messages, err := m.service.GetMessagesByChatID(friendID, UserID)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, "Ошибка при отправки данных сообщения", http.StatusInternalServerError)
	}
}
