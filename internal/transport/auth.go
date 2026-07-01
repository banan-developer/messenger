package transport

import (
	"fmt"
	"messenger_v2/internal/service"
	"messenger_v2/pkg/auth"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./web/html/login.html")
		return
	}

	if r.Method == http.MethodPost {
		login := r.FormValue("email")
		password := r.FormValue("password")
		UserID, err := a.authService.Login(login, password)

		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		auth.SetUserID(w, r, UserID)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)

	}
}

func (a *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./web/html/registration.html")
		return
	}
	login := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")
	sex := r.FormValue("sex")

	err := a.authService.Registration(login, password, name, sex)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Пользователь с таким email уже существует", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessions(w, r)
	w.WriteHeader(http.StatusOK)
}
