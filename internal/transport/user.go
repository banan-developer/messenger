package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/service"
	"messenger_v2/pkg/auth"
	"net/http"
	"os"
	"strconv"
	"time"
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
		idSTR := r.URL.Query().Get("id")
		if idSTR != "" {
			h.GetPersonByID(w, r)
		} else {
			h.GetProfile(w, r)
		}
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
	userID, _ := auth.GetUserId(r)

	user, err := h.service.GetProfile(userID)

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

func (h *UserHandler) UploadAvatarUser(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("avatar")

	if err != nil {
		http.Error(w, "Ошибка загрузки", 500)
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	path := "./web/static/uploads/avatars/" + filename

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Ошибка сохранения", 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	avatarURL := "/static/uploads/avatars/" + filename

	err = h.service.UploadAvatarUser(UserID, avatarURL)
	if err != nil {
		log.Println("Ошибка при вызове сервиса UploadAvatarUser")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"avatar": avatarURL,
	})
}

func (h *UserHandler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("id")
	PersonID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	person, err := h.service.GetPersonByID(PersonID)
	if err != nil {
		fmt.Println("REAL DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(person)
	if err != nil {
		http.Error(w, "Ошибка при отравки данных о пользователе", http.StatusInternalServerError)
	}

}
