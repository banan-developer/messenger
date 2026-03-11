package main

import (
	"messenger/auth"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (app *application) autoresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./pkg/ui/html/auth.html")
		return
	}
	// получения значения в input через поля
	if r.Method == http.MethodPost {
		login := r.FormValue("email")
		password := r.FormValue("password")

		var UserId int
		var PasswordFromBd string

		// получаем данные из бд, а потом сравниваем пароль из базы данных и написанным в input
		rows, err := app.db.Query("SELECT id, login, password FROM users WHERE login = ?", login)

		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		err = rows.Scan(&UserId, &PasswordFromBd)

		// проверка пароля на совпадение, PasswordFromdb хэшированные, password полученый из поля input
		HashError := bcrypt.CompareHashAndPassword(
			[]byte(PasswordFromBd),
			[]byte(password),
		)
		if HashError == nil {
			auth.SetUserID(w, r, UserId)
		} else {
			return // вот тут добавить не забудь, логиравание ошибок
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) regHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./pkg/ui/html/registration.html")
		return
	}

	// получение значения в input через поля
	login := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

	hashedPassword, _ := hashPassword(password)
	_, err := app.db.Exec("INSERT INTO users (login, password, name) VALUES (?, ?, ?)", login, hashedPassword, name)
	if err != nil {
		app.errorLog.Println("REGISTER ERROR:", err)
		http.Error(w, "Пользователь с таким email уже существует", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) exitSession(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessions(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
