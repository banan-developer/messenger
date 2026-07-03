package domain

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	About    string `json:"about"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
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
