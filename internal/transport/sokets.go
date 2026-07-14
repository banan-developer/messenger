package transport

import (
	"encoding/json"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/service"
	"messenger_v2/internal/transport/websocket"
	"net/http"
	"strconv"
	"time"
)

type WebSocketHandler struct {
	msgService *service.MessageService
	hub        *websocket.Hub
}

func NewWebSocketHandler(msgService *service.MessageService, hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		msgService: msgService,
		hub:        hub,
	}
}

// HandleWS — WebSocket соединение
func (h *WebSocketHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	if chatID, err := strconv.Atoi(r.URL.Query().Get("chat_id")); err == nil && chatID > 0 {
		userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
		h.hub.HandleWebSocket(w, r, userID, chatID, func(data []byte) { h.ProcessMessage(userID, chatID, data) })
		return
	}
	friendID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || friendID <= 0 {
		http.Error(w, "Invalid friend id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	log.Printf("WebSocket connection: user=%d, friend=%d", userID, friendID)

	// Получаем или создаём чат
	chatID, err := h.msgService.GetOrCreateChat(userID, friendID)
	if err != nil {
		log.Printf("Failed to get/create chat: %v", err)
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	log.Printf("WebSocket connected: user=%d, friend=%d, chat=%d", userID, friendID, chatID)

	// Обрабатываем WebSocket
	h.hub.HandleWebSocket(w, r, userID, chatID, func(data []byte) {
		h.ProcessMessage(userID, chatID, data)
	})
}

// ProcessMessage — обработка входящего сообщения из WebSocket
func (h *WebSocketHandler) ProcessMessage(userID, chatID int, data []byte) {
	var req struct {
		Event string `json:"event"`
		ID    int    `json:"id"`
		Text  string `json:"text"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("Invalid WebSocket message from user %d: %v", userID, err)
		return
	}

	switch req.Event {
	case "edit_message":
		h.editMessage(userID, chatID, req.ID, req.Text)
		return
	case "delete_message":
		h.deleteMessage(userID, chatID, req.ID)
		return
	case "":
		// Обычное новое сообщение.
	default:
		log.Printf("Unknown WebSocket event from user %d: %s", userID, req.Event)
		return
	}

	if req.Text == "" {
		log.Printf("Empty message from user %d", userID)
		return
	}

	log.Printf("Processing message from user %d, chat %d: %s", userID, chatID, req.Text)

	// Сохраняем сообщение
	msg := &domain.Message{
		Text:   req.Text,
		FromID: userID,
		ChatID: chatID,
	}

	err := h.msgService.SaveMessage(msg)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	log.Printf("Message saved: id=%d, chat=%d", msg.ID, chatID)

	// Отправляем сообщение всем участникам чата
	response := map[string]interface{}{
		"event": "new_message",
		"message": map[string]interface{}{
			"id":         msg.ID,
			"text":       msg.Text,
			"from_id":    msg.FromID,
			"to_id":      msg.ToID,
			"chat_id":    msg.ChatID,
			"created_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	h.hub.SendToChat(chatID, response, 0)
}

func (h *WebSocketHandler) editMessage(userID, chatID, msgID int, text string) {
	if msgID <= 0 || text == "" {
		return
	}
	msg, err := h.msgService.GetMessageByID(msgID)
	if err != nil || msg.ChatID != chatID || msg.FromID != userID {
		log.Printf("Rejected message edit: user=%d, message=%d", userID, msgID)
		return
	}
	if err := h.msgService.UpdateMessage(msgID, text); err != nil {
		log.Printf("Failed to edit message %d: %v", msgID, err)
		return
	}
	h.hub.SendToChat(chatID, map[string]interface{}{
		"event": "message_updated",
		"data": map[string]interface{}{
			"id":   msgID,
			"text": text,
		},
	}, 0)
}

func (h *WebSocketHandler) deleteMessage(userID, chatID, msgID int) {
	if msgID <= 0 {
		return
	}
	msg, err := h.msgService.GetMessageByID(msgID)
	if err != nil || msg.ChatID != chatID || msg.FromID != userID {
		log.Printf("Rejected message delete: user=%d, message=%d", userID, msgID)
		return
	}
	if err := h.msgService.DeleteMessage(msgID); err != nil {
		log.Printf("Failed to delete message %d: %v", msgID, err)
		return
	}
	h.hub.SendToChat(chatID, map[string]interface{}{
		"event": "message_deleted",
		"data": map[string]interface{}{
			"id": msgID,
		},
	}, 0)
}

// ProcessFileMessage — обработка сообщения с файлом из WebSocket
func (h *WebSocketHandler) ProcessFileMessage(userID, chatID int, data []byte) {
	var req struct {
		Text          string `json:"text"`
		AttachmentURL string `json:"attachment_url"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("Invalid WebSocket file message from user %d: %v", userID, err)
		return
	}

	if req.AttachmentURL == "" {
		log.Printf("Empty attachment from user %d", userID)
		return
	}

	log.Printf("Processing file message from user %d, chat %d: %s", userID, chatID, req.AttachmentURL)

	// Сохраняем сообщение
	msg := &domain.Message{
		Text:          req.Text,
		FromID:        userID,
		ChatID:        chatID,
		AttachmentURL: req.AttachmentURL,
	}

	err := h.msgService.SaveMessage(msg)
	if err != nil {
		log.Printf("Failed to save file message: %v", err)
		return
	}

	// Отправляем сообщение всем участникам чата
	response := map[string]interface{}{
		"event": "new_message",
		"message": map[string]interface{}{
			"id":             msg.ID,
			"text":           msg.Text,
			"from_id":        msg.FromID,
			"to_id":          msg.ToID,
			"chat_id":        msg.ChatID,
			"attachment_url": msg.AttachmentURL,
			"created_at":     time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	h.hub.SendToChat(chatID, response, 0)
}
