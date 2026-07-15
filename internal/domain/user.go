package domain

// User описывает полный профиль пользователя в системе.
// Превращается в JSON при запросе собственного профиля (GET /api/profile) или при его обновлении (PUT /api/profile).
type User struct {
	// ID уникальный идентификатор пользователя.
	ID       int    `json:"id"`
	Name     string `json:"name"`
	About    string `json:"about"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
}

type UserResponse struct {
	Name   string `json:"name"`
	About  string `json:"about"`
	Avatar string `json:"avatar"`
	Sex    string `json:"sex"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegistrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
}
