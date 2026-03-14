package main

import (
	"encoding/json"
	"fmt"
	"io"
	"messenger/auth"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./pkg/ui/html/profile.html")
}

func (app *application) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // переход из http в websocket
	if err != nil {
		fmt.Println("Ошибка при переходе на Websocket протокол")
		return
	}
	client[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("client disconnected")
			delete(client, conn)
			conn.Close()
			break
		}
		broadcast <- message

	}
}

// Работа с данными профиля
func (app *application) getPerson(w http.ResponseWriter, r *http.Request) {
	var person person

	UserID, _ := auth.GetUserId(r)

	err := app.db.QueryRow("SELECT name, about, avatar_url, sex FROM users WHERE id = ?", UserID).Scan(&person.Name, &person.About, &person.Avatar, &person.Sex)
	if err != nil {
		fmt.Println("REAL DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(person)
	if err != nil {
		fmt.Println("error2!")
	}

}

func (app *application) updateProfile(w http.ResponseWriter, r *http.Request) {
	var person person
	UserID, _ := auth.GetUserId(r)

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Ошибка при получения данных с фронта", http.StatusInternalServerError)
		return
	}
	_, err = app.db.Exec("UPDATE users SET name = ?, about = ? WHERE id = ?", person.Name, person.About, UserID)
	if err != nil {
		http.Error(w, "Ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}
}

// Занесение поля в бд
func (app *application) GetWall(w http.ResponseWriter, r *http.Request) {
	var wall wall
	UserID, _ := auth.GetUserId(r)
	err := json.NewDecoder(r.Body).Decode(&wall)
	if err != nil {
		http.Error(w, "Ошибка при получении данных с фронта", http.StatusInternalServerError)
		return
	}

	_, err = app.db.Exec("INSERT INTO wall (users_id, title, text) VALUES (?, ?, ?)", UserID, wall.Title, wall.Text)
	if err != nil {
		app.errorLog.Println("Ошибка в базе данных", err)
	}

}

func (app *application) PushwWall(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)
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

// удаление поста
func (app *application) deleteWall(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idSTR)
	if err != nil {
		app.errorLog.Println("DB DELETE ERROR:", err)
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	UserID, _ := auth.GetUserId(r)

	_, err = app.db.Exec("DELETE FROM wall WHERE idwall = ? AND users_id = ?", id, UserID)

	if err != nil {
		app.errorLog.Println("DB DELETE ERROR:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// редактирования поста

func (app *application) editWall(w http.ResponseWriter, r *http.Request) {
	var wall wall
	idSTR := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idSTR)
	if err != nil {
		app.errorLog.Println("DB DELETE ERROR:", err)
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	err = app.db.QueryRow("SELECT title FROM wall WHERE idwall = ?", id).Scan(&wall.Title)
	if err != nil {
		app.errorLog.Println("Note not found:", err)
		return
	}
	err = app.db.QueryRow("SELECT text FROM wall WHERE idwall = ?", id).Scan(&wall.Text)
	if err != nil {
		app.errorLog.Println("Note not found:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(wall)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func (app *application) updateEditingWall(w http.ResponseWriter, r *http.Request) {
	var wall wall

	err := json.NewDecoder(r.Body).Decode(&wall)
	idSTR := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idSTR)
	if err != nil {
		app.errorLog.Println("DB DELETE ERROR:", err)
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}

	_, err = app.db.Exec("UPDATE wall SET title = ?, text = ? WHERE idwall = ?", wall.Title, wall.Text, id)
}

func (app *application) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	UserID, _ := auth.GetUserId(r)
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("avatar")

	if err != nil {
		http.Error(w, "Ошибка загрузки", 500)
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	path := "./pkg/ui/static/avatars/" + filename

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Ошибка сохранения", 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	avatarURL := "/static/avatars/" + filename

	_, err = app.db.Exec(
		"UPDATE users SET avatar_url = ? WHERE id = ? ",
		avatarURL,
		UserID,
	)

	json.NewEncoder(w).Encode(map[string]string{
		"avatar": avatarURL,
	})
}

func (app *application) foundUsers(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	rows, err := app.db.Query("SELECT id, name, avatar_url FROM users WHERE name = ?", name)
	if err != nil {
		http.Error(w, "Ошибка при получении имени пользователя для его поиска", 500)
	}

	defer rows.Close()
	var per []person

	for rows.Next() {
		var p person
		rows.Scan(&p.ID, &p.Name, &p.Avatar)
		per = append(per, p)
	}

	if per == nil {
		per = []person{}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(per)
	if err != nil {
		http.Error(w, "Ошибка при отправки данных пользователя при запросе добавить в друзья", http.StatusInternalServerError)
	}

}

func (app *application) AddToFriend(w http.ResponseWriter, r *http.Request) {
	UserId, _ := auth.GetUserId(r)
	idSTR := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idSTR)
	if err != nil {
		app.errorLog.Println("DB DELETE ERROR:", err)
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	status := "accepted"

	_, err = app.db.Exec("INSERT INTO friends (friend_id, users_id, status) VALUES (?, ?, ?)", id, UserId, status)
	if err != nil {
		app.errorLog.Println("DB QUERY ERROR", err)
		return
	}

}

func (app *application) loadFriend(w http.ResponseWriter, r *http.Request) {
	UserId, _ := auth.GetUserId(r)

	rows, err := app.db.Query(
		`SELECT users.id, users.name, users.avatar_url
		FROM friends
		JOIN users ON users.id = friends.friend_id
		WHERE friends.users_id = ?
		AND friends.status = 'accepted';
		`, UserId)

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
