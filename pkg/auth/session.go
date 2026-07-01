package auth

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

// создание сессии
var Store = sessions.NewCookieStore([]byte("very-secret-super-long-key-1234567890"))

func InitStore() {
	gob.Register(int(0))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}
}

const sessionName = "notes-session"
const userIDKey = "user_id"

// кидаем айди пользователя в куки
func SetUserID(w http.ResponseWriter, r *http.Request, userID int) error {
	sessions, _ := Store.Get(r, sessionName)
	sessions.Values[userIDKey] = userID
	return sessions.Save(r, w)

}

// получем айди
func GetUserId(r *http.Request) (int, bool) {
	sessions, err := Store.Get(r, sessionName)
	if err != nil {
		return 0, false
	}
	id, ok := sessions.Values[userIDKey].(int)
	return id, ok
}

// выход из сессии
func ClearSessions(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1

	return sessions.Save(r, w)
}
