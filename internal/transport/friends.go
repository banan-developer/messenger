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
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func (f *FriendHandler) GetFriendByID(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)

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

	f.service.AddToFriend(UserID, FriendID, status)
}

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
