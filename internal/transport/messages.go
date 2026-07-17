package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/service"
	"messenger_v2/internal/transport/websocket"
	"messenger_v2/pkg/auth"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MessageHandler struct {
	msgService *service.MessageService
	hub        *websocket.Hub
}

func NewMessageHandler(msgService *service.MessageService, hub *websocket.Hub) *MessageHandler {
	return &MessageHandler{
		msgService: msgService,
		hub:        hub,
	}
}

// Messages возвращает список чатов текущего пользователя либо сообщения конкретного чата.
//
// # API Контракт
//
//	Маршрут:     /api/messages
//	Авторизация: Требуется (сессионная кука)
//
//	GET /api/messages                  — получить список чатов текущего пользователя
//	GET /api/messages?chat_id={id}     — получить сообщения конкретного чата
//	GET /api/messages?id={friend_id}   — открыть/создать чат с другом и получить его сообщения
//
// # Параметры запроса
//
//	chat_id (int, необяз.) — id чата с сообщениями
//	id      (int, необяз.) — id друга для открытия личного чата
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON со списком чатов или сообщениями
func (h *MessageHandler) Messages(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("id") == "" && r.URL.Query().Get("chat_id") == "" {
		userID, _ := auth.GetUserId(r)
		chats, err := h.msgService.GetChatsForUser(userID)
		if err != nil {
			http.Error(w, "Failed to get chats", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chats)
		return
	}
	if chatID, err := strconv.Atoi(r.URL.Query().Get("chat_id")); err == nil && chatID > 0 {
		userID, _ := auth.GetUserId(r)
		allowed, err := h.msgService.IsChatParticipant(chatID, userID)
		if err != nil || !allowed {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		messages, err := h.msgService.GetMessagesByChatID(chatID)
		if err != nil {
			http.Error(w, "Failed to get messages", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
		return
	}
	friendID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || friendID <= 0 {
		http.Error(w, "Invalid friend id", http.StatusBadRequest)
		return
	}

	UserID, _ := auth.GetUserId(r)

	chatID, err := h.msgService.GetOrCreateChat(UserID, friendID)
	if err != nil {
		log.Println("GetOrCreateChat error:", err)
		http.Error(w, "Failed to get chat", http.StatusInternalServerError)
		return
	}

	messages, err := h.msgService.GetMessagesByChatID(chatID)
	if err != nil {
		log.Println("GetMessagesByChatID error:", err)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetMessageByID возвращает сообщение по его идентификатору.
//
// # API Контракт
//
//	Метод:       GET
//	Маршрут:     /api/messages
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	msg_id (int, обяз.) — id сообщения
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON с данными сообщения
func (h *MessageHandler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	msgID, err := strconv.Atoi(r.URL.Query().Get("msg_id"))
	if err != nil || msgID <= 0 {
		http.Error(w, "Invalid msg_id", http.StatusBadRequest)
		return
	}

	msg, err := h.msgService.GetMessageByID(msgID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// UpdateMessage редактирует текст существующего сообщения.
//
// # API Контракт
//
//	Метод:       PUT
//	Маршрут:     /api/messages
//	Авторизация: Требуется (сессионная кука)
//
// # Формат данных запроса
//
//	Content-Type: application/json
//	Тело: JSON-объект с полями id и text
func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	var req domain.EditMessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID <= 0 {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	msg, _ := h.msgService.GetMessageByID(req.ID)

	err = h.msgService.UpdateMessage(req.ID, req.Text)
	if err != nil {
		log.Println("UpdateMessage error:", err)
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}

	if msg != nil {
		h.hub.SendToChat(msg.ChatID, map[string]interface{}{
			"event": "message_updated",
			"data": map[string]interface{}{
				"id":   req.ID,
				"text": req.Text,
			},
		}, 0)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DeleteMessage удаляет сообщение по его идентификатору.
//
// # API Контракт
//
//	Метод:       DELETE
//	Маршрут:     /api/messages
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	id (int, обяз.) — id сообщения для удаления
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	msgID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || msgID <= 0 {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	msg, _ := h.msgService.GetMessageByID(msgID)

	err = h.msgService.DeleteMessage(msgID)
	if err != nil {
		log.Println("DeleteMessage error:", err)
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	if msg != nil {
		h.hub.SendToChat(msg.ChatID, map[string]interface{}{
			"event": "message_deleted",
			"data": map[string]interface{}{
				"id": msgID,
			},
		}, 0)
	}

	w.WriteHeader(http.StatusNoContent)
}

// SendMessageWithFile отправляет сообщение с прикрепленным файлом в чат или личный диалог.
//
// # API Контракт
//
//	Метод:       POST
//	Маршрут:     /api/messages/file
//	Авторизация: Требуется (сессионная кука)
//
// # Формат данных запроса
//
//	Формат: multipart/form-data
//	chat_id   (int, необяз.) — id чата, если сообщение отправляется в групповой чат
//	friend_id (int, необяз.) — id друга, если нужно открыть или создать личный чат
//	text      (string, необяз.) — текст сообщения
//	file      (file, обяз.) — файл вложения
func (h *MessageHandler) SendMessageWithFile(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)
	if UserID <= 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	const maxFileSize = 20 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize+(1<<20))
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(r.FormValue("chat_id"))
	if chatID > 0 {
		allowed, accessErr := h.msgService.IsChatParticipant(chatID, UserID)
		if accessErr != nil || !allowed {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	} else {
		friendID, friendErr := strconv.Atoi(r.FormValue("friend_id"))
		if friendErr != nil || friendID <= 0 {
			http.Error(w, "Invalid friend_id", http.StatusBadRequest)
			return
		}
		chatID, err = h.msgService.GetOrCreateChat(UserID, friendID)
		if err != nil {
			http.Error(w, "Failed to get chat", http.StatusInternalServerError)
			return
		}
	}

	text := r.FormValue("text")
	var attachmentURL string

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()
	if handler.Size <= 0 || handler.Size > maxFileSize {
		http.Error(w, "File is empty or exceeds 20 MB", http.StatusRequestEntityTooLarge)
		return
	}

	ext := strings.ToLower(filepath.Ext(filepath.Base(handler.Filename)))
	allowedExtensions := map[string]string{
		".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".gif": "image/gif", ".webp": "image/webp", ".pdf": "application/pdf",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xml":  "application/xml",
	}
	attachmentType, allowed := allowedExtensions[ext]
	if !allowed {
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return
	}

	headerBytes := make([]byte, 512)
	readCount, readErr := file.Read(headerBytes)
	if readErr != nil && readErr != io.EOF {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	detectedType := http.DetectContentType(headerBytes[:readCount])
	validContent := strings.HasPrefix(attachmentType, "image/") && strings.HasPrefix(detectedType, "image/")
	switch ext {
	case ".pdf":
		validContent = detectedType == "application/pdf"
	case ".docx":
		validContent = detectedType == "application/zip" || detectedType == "application/octet-stream"
	case ".xml":
		validContent = detectedType == "text/xml; charset=utf-8" ||
			detectedType == "text/plain; charset=utf-8" ||
			detectedType == "application/octet-stream"
	}
	if !validContent {
		http.Error(w, "File content does not match its extension", http.StatusUnsupportedMediaType)
		return
	}

	const uploadsDir = "./web/static/uploads/files"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Println("Upload directory error:", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(uploadsDir, filename)
	dst, err := os.Create(path)
	if err != nil {
		log.Println("File save error:", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		log.Println("File copy error:", err)
		os.Remove(path)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	if err := dst.Close(); err != nil {
		os.Remove(path)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	attachmentURL = "/static/uploads/files/" + filename

	msg := &domain.Message{
		Text:           text,
		FromID:         UserID,
		ChatID:         chatID,
		AttachmentURL:  attachmentURL,
		AttachmentName: filepath.Base(handler.Filename),
		AttachmentType: attachmentType,
		AttachmentSize: handler.Size,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}

	err = h.msgService.SaveMessage(msg)
	if err != nil {
		log.Println("Save message error:", err)
		os.Remove(path)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	h.hub.SendToChat(chatID, map[string]interface{}{
		"event": "new_message",
		"message": map[string]interface{}{
			"id":              msg.ID,
			"text":            msg.Text,
			"from_id":         msg.FromID,
			"chat_id":         msg.ChatID,
			"created_at":      msg.CreatedAt,
			"attachment_url":  msg.AttachmentURL,
			"attachment_name": msg.AttachmentName,
			"attachment_type": msg.AttachmentType,
			"attachment_size": msg.AttachmentSize,
		},
	}, 0)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":              msg.ID,
		"text":            msg.Text,
		"attachment_url":  msg.AttachmentURL,
		"attachment_name": msg.AttachmentName,
		"attachment_type": msg.AttachmentType,
		"attachment_size": msg.AttachmentSize,
		"created_at":      msg.CreatedAt,
	})
}
