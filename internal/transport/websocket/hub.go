package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients map[int]*client
	mu      sync.RWMutex
}

type client struct {
	conn   *websocket.Conn
	chatID int
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]*client),
	}
}

// HandleWebSocket — обрабатывает WebSocket соединение
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID, chatID int, onMessage func([]byte)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error for user %d: %v", userID, err)
		return
	}

	h.mu.Lock()
	h.clients[userID] = &client{conn: conn, chatID: chatID}
	h.mu.Unlock()

	log.Printf("User %d connected to WebSocket", userID)

	defer func() {
		h.mu.Lock()
		delete(h.clients, userID)
		h.mu.Unlock()
		conn.Close()
		log.Printf("User %d disconnected from WebSocket", userID)
	}()

	// Передаём полученные сообщения обработчику, который сохраняет их в БД.
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %d: %v", userID, err)
			} else {
				log.Printf("WebSocket closed for user %d: %v", userID, err)
			}
			break
		}
		onMessage(data)
	}
}

// SendToUser — отправить сообщение конкретному пользователю
func (h *Hub) SendToUser(userID int, data interface{}) error {
	h.mu.RLock()
	conn, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return nil
	}
	return conn.conn.WriteJSON(data)
}

// SendToChat — отправить сообщение всем участникам чата (кроме excludeUserID)
func (h *Hub) SendToChat(chatID int, data interface{}, excludeUserID int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, client := range h.clients {
		if userID == excludeUserID {
			continue
		}
		if client.chatID != chatID {
			continue
		}
		if err := client.conn.WriteJSON(data); err != nil {
			log.Printf("Error sending to user %d: %v", userID, err)
		}
	}
}

// SendToAll — отправить сообщение всем подключённым пользователям
func (h *Hub) SendToAll(data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, client := range h.clients {
		if err := client.conn.WriteJSON(data); err != nil {
			log.Printf("Error sending to user %d: %v", userID, err)
		}
	}
}

// GetClients — получить список всех подключённых клиентов
func (h *Hub) GetClients() map[int]*client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients
}

// IsOnline — проверить, онлайн ли пользователь
func (h *Hub) IsOnline(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
