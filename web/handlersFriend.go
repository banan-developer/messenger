package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) getAnotherProfile(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("id")
	profileID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	var person person
	err = app.db.QueryRow("SELECT name, about, avatar_url, sex FROM users WHERE id = ?", profileID).Scan(&person.Name, &person.About, &person.Avatar, &person.Sex)
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
func (app *application) getWallFromAnotherProfile(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("user_id")
	UserID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	rows, err := app.db.Query("SELECT idwall, title, text FROM wall WHERE users_id  = ?", UserID)
	if err != nil {
		app.errorLog.Println("DB QUERY ERROR:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var content []wall

	for rows.Next() {
		var wall wall
		rows.Scan(&wall.Id, &wall.Title, &wall.Text)
		content = append(content, wall)
	}
	if content == nil {
		content = []wall{}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(content)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (app *application) loadAnotherFriend(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("user_id")
	UserID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	rows, err := app.db.Query(
		`SELECT users.id, users.name, users.avatar_url
		FROM friends
		JOIN users ON users.id = friends.friend_id
		WHERE friends.users_id = ?
		AND friends.status = 'accepted';
		`, UserID)

	if err != nil {
		http.Error(w, "Ошибка при получении имени пользователя для его поиска", 500)
	}

	defer rows.Close()

	var user []person
	for rows.Next() {
		var person person
		rows.Scan(&person.ID, &person.Name, &person.Avatar)
		user = append(user, person)
	}
	if user == nil {
		user = []person{}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Ошибка при отправки данных пользователя при запросе добавить в друзья", http.StatusInternalServerError)
	}
}
