package main

import "net/http"

func (app *application) handleProfile(w http.ResponseWriter, r *http.Request) {
	profileID := r.URL.Query().Get("id")
	switch r.Method {
	case http.MethodGet:
		if profileID != "" {
			app.getAnotherProfile(w, r)
		} else {
			app.getPerson(w, r)
		}
	case http.MethodPut:
		app.updateProfile(w, r)
	default:
		http.Error(w, "Метод не найден", http.StatusMethodNotAllowed)
	}
}

func (app *application) handleProfileAvatar(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.uploadAvatar(w, r)
	default:
		http.Error(w, "Метод не найден", http.StatusMethodNotAllowed)
	}
}

func (app *application) handlePost(w http.ResponseWriter, r *http.Request) {
	PostID := r.URL.Query().Get("id")
	UserID := r.URL.Query().Get("user_id")
	switch r.Method {
	case http.MethodGet:
		if PostID != "" {
			app.editWall(w, r)
		} else if UserID != "" {
			app.getWallFromAnotherProfile(w, r)
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
	userID := r.URL.Query().Get("user_id")
	switch r.Method {
	case http.MethodGet:
		if searchName != "" {
			app.foundUsers(w, r)
		} else if userID != "" {
			app.loadAnotherFriend(w, r)
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
