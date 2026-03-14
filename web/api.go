package main

import "net/http"

func (app *application) handleProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getPerson(w, r)
	case http.MethodPut:
		app.updateProfile(w, r)
	case http.MethodPost:
		app.uploadAvatar(w, r)
	default:
		http.Error(w, "Метод не найден", http.StatusMethodNotAllowed)
	}
}

func (app *application) handlePost(w http.ResponseWriter, r *http.Request) {
	PostID := r.URL.Query().Get("id")
	switch r.Method {
	case http.MethodGet:
		if PostID != "" {
			app.editWall(w, r)
		} else {
			app.PushwWall(w, r)
		}
	case http.MethodPost:
		app.GetWall(w, r)
	case http.MethodPut:
		if PostID != "" {
			app.updateEditingWall(w, r)
		} else {
			http.Error(w, "ID поста не указан", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if PostID != "" {
			app.deleteWall(w, r)
		} else {
			http.Error(w, "ID поста не указан", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Метод не найден", http.StatusMethodNotAllowed)
	}
}

func (app *application) handleFriend(w http.ResponseWriter, r *http.Request) {
	friendID := r.URL.Query().Get("id")
	searchName := r.URL.Query().Get("name")
	switch r.Method {
	case http.MethodGet:
		if searchName != "" {
			app.foundUsers(w, r)
		} else {
			app.loadFriend(w, r)
		}
	case http.MethodPost:
		if friendID != "" {
			app.AddToFriend(w, r)
		} else {
			http.Error(w, "ID друга не указан", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Метод не найден", http.StatusMethodNotAllowed)
	}
}
