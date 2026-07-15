package transport

import (
	"fmt"
	"messenger_v2/internal/domain"
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

// Login обрабатывает вход пользователя в систему.
//
// # API Контракт
//
//	Метод:       POST
//	Маршрут:     /login
//	Авторизация: Не требуется
//
// # Параметры запроса
//
//	Формат: application/x-www-form-urlencoded
//	email    (string, обяз.) — email или логин пользователя
//	password (string, обяз.) — пароль пользователя
//
// # Ответы
//
//	303 See Other    — успешный вход, установка сессионной куки и редирект на /profile
//	401 Unauthorized — неверные учетные данные
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

// Registration обрабатывает создание новой учетной записи пользователя.
//
// # API Контракт
//
//	Метод:       POST
//	Маршрут:     /registration
//	Авторизация: Не требуется
//
// # Параметры запроса
//
// Формат: application/x-www-form-urlencoded
//
//	email    (string, обяз.) — Электронная почта
//	password (string, обяз.) — Пароль (от 6 символов)
//	name     (string, обяз.) — Имя или ФИО (до 100 символов)
//	sex      (string, обяз.) — Пол ("Мужской" или "Женский")
//
// # Ответы
//
//	303 See Other             — успешная регистрация, редирект на /login
//	500 Internal Server Error — email уже занят или ошибка БД
func (a *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./web/html/registration.html")
		return
	}
	res := &domain.RegistrationRequest{
		Login:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Sex:      r.FormValue("sex"),
	}

	err := a.authService.Registration(res)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Пользователь с таким email уже существует", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Logout завершает текущую сессию пользователя и очищает авторизационные куки.
//
// # API Контракт
//
//	Метод:       GET
//	Маршрут:     /exit
//	Авторизация: Требуется (сессионная кука)
//
// # Ответы
//
//	200 OK — сессия успешно завершена
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessions(w, r)
	w.WriteHeader(http.StatusOK)
}
