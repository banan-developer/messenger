package transport

import (
	"encoding/json"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/service"
	"messenger_v2/pkg/auth"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProfile(w, r)
	case http.MethodPut:
		h.UpdateProfile(w, r)
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	UserID, _ := auth.GetUserId(r)

	user, err := h.service.GetProfile(UserID)

	if err != nil {
		log.Printf("GetProfile error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	UserID, _ := auth.GetUserId(r)
	user.ID = UserID
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Ошибка при получения данных с фронта", http.StatusInternalServerError)
		return
	}

	err = h.service.UpdateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
