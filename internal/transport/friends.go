package transport

import (
	"encoding/json"
	"log"
	"messenger_v2/internal/service"
	"messenger_v2/pkg/auth"
	"net/http"
	"strconv"
)

type FriendHandler struct {
	service *service.FriendService
}

func NewFriendHandler(service *service.FriendService) *FriendHandler {
	return &FriendHandler{
		service: service,
	}
}

// Friends является маршрутизатором для управления друзьями и заявками.
//
// # API Контракт
//
//	Маршрут:     /api/friend
//	Авторизация: Требуется (сессионная кука)
//
//	GET  /api/friend?name={name}      — поиск пользователей по имени
//	GET  /api/friend?user_id={id}    — получение списка друзей пользователя
//	GET  /api/friend                 — получение списка друзей текущего пользователя
//	POST /api/friend?id={id}         — отправка заявки в друзья
//	PUT  /api/friend?friendID={id}   — принятие входящей заявки
//	DELETE /api/friend?id={id}       — удаление пользователя из друзей
//
// # Параметры запроса
//
//	name     (string, необяз.) — строка для поиска пользователей
//	user_id  (int, необяз.)    — id пользователя, чьи друзья запрашиваются
//	id       (int, обяз. для POST/DELETE) — id целевого пользователя
//	friendID (int, обяз. для PUT)         — id пользователя, заявку которого принимают
//
// # Ответы
//
//	200 OK                — успешная обработка запроса
//	204 No Content        — успешное удаление из друзей
//	400 Bad Request       — неверный формат параметров URL
//	405 Method Not Allowed — неподдерживаемый HTTP-метод
func (f *FriendHandler) Friends(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		FriendName := r.URL.Query().Get("name")
		if FriendName != "" {
			f.FoundFriendByID(w, r)
		} else {
			f.GetFriendByID(w, r)
		}
	case http.MethodPost:
		f.AddToFriend(w, r)
	case http.MethodPut:
		f.AcceptComingRequset(w, r)
	case http.MethodDelete:
		userID, _ := auth.GetUserId(r)
		friendID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid friend id", 400)
			return
		}
		if err = f.service.DeleteFriend(userID, friendID); err != nil {
			http.Error(w, "Failed", 500)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

// GetFriendByID возвращает список друзей указанного пользователя или текущего авторизованного пользователя.
//
// # API Контракт
//
//	Метод:       GET
//	Маршрут:     /api/friend
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	user_id (int, необяз.) — id пользователя, чьи друзья нужно получить.
//	                         Если не передан, используется текущий пользователь из сессии.
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON-массив объектов друзей
func (f *FriendHandler) GetFriendByID(w http.ResponseWriter, r *http.Request) {
	UseridSTR := r.URL.Query().Get("user_id")
	var UserID int
	var err error
	if UseridSTR != "" {
		UserID, err = strconv.Atoi(UseridSTR)
		if err != nil {
			http.Error(w, "Invalid note id", http.StatusBadRequest)
			return
		}
	} else {
		UserID, _ = auth.GetUserId(r)
	}

	friends, err := f.service.GetFriendsByID(UserID)
	if err != nil {
		http.Error(w, "Ошибка при получении имени пользователя для его поиска", 500)
	}

	err = json.NewEncoder(w).Encode(friends)
	if err != nil {
		http.Error(w, "Ошибка при отправки данных пользователя при запросе добавить в друзья", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}

// AddToFriend отправляет заявку в друзья другому пользователю.
//
// # API Контракт
//
//	Метод:       POST
//	Маршрут:     /api/friend
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	id (int, обяз.) — id пользователя, которому отправляется заявка
//
// # Формат запроса
//
//	Тело запроса отсутствует. Параметр передается через query string.
func (f *FriendHandler) AddToFriend(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)

	idSTR := r.URL.Query().Get("id")

	FriendID, err := strconv.Atoi(idSTR)
	if err != nil {
		log.Println("DB ERROR:", err)
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	status := "invited"

	if err := f.service.AddToFriend(UserID, FriendID, status); err != nil {
		http.Error(w, "Failed to send friend request", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// FoundFriendByID ищет пользователей по имени или части имени.
//
// # API Контракт
//
//	Метод:       GET
//	Маршрут:     /api/friend
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	name (string, обяз.) — подстрока для поиска пользователя по имени
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON-массив найденных пользователей
func (f *FriendHandler) FoundFriendByID(w http.ResponseWriter, r *http.Request) {
	FriendName := r.URL.Query().Get("name")
	friend, err := f.service.FoundFriendByID(FriendName)
	if err != nil {
		http.Error(w, "Ошибка при получении данных", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(friend)
	if err != nil {
		http.Error(w, "Ошибка при отправки данных пользователя при запросе добавить в друзья", http.StatusInternalServerError)
	}
}

// GetIncomigRequest возвращает список входящих заявок в друзья для текущего пользователя.
//
// # API Контракт
//
//	Метод:       GET
//	Маршрут:     /api/incomingrequest
//	Авторизация: Требуется (сессионная кука)
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON-массив заявок в друзья
func (f *FriendHandler) GetIncomigRequest(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)

	friendRequest, err := f.service.GetIncomingRequest(UserID)
	if err != nil {
		http.Error(w, "Ошибка получения данных", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(friendRequest)
	if err != nil {
		http.Error(w, "Ошибка при отправки данных входящих заявок", http.StatusInternalServerError)
	}
}

// AcceptComingRequset принимает входящую заявку в друзья.
//
// # API Контракт
//
//	Метод:       PUT
//	Маршрут:     /api/friend
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	friendID (int, обяз.) — id пользователя, чью заявку принимают
//
// # Формат запроса
//
//	Тело запроса отсутствует. ID передается через query string.
func (f *FriendHandler) OutgoingRequests(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserId(r)

	switch r.Method {
	case http.MethodGet:
		requests, err := f.service.GetOutgoingRequests(userID)
		if err != nil {
			http.Error(w, "Failed to get outgoing requests", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(requests); err != nil {
			http.Error(w, "Failed to encode outgoing requests", http.StatusInternalServerError)
		}
	case http.MethodDelete:
		friendID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || friendID <= 0 {
			http.Error(w, "Invalid friend id", http.StatusBadRequest)
			return
		}
		if err := f.service.CancelOutgoingRequest(userID, friendID); err != nil {
			http.Error(w, "Failed to cancel request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func (f *FriendHandler) AcceptComingRequset(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)

	FriendidSTR := r.URL.Query().Get("friendID")
	FriendID, err := strconv.Atoi(FriendidSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	err = f.service.AcceptComingRequset(FriendID, UserID)
	if err != nil {
		http.Error(w, "Ошибка при обнавлении данных", http.StatusInternalServerError)
	}

}
