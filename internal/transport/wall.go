package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/service"
	"messenger_v2/pkg/auth"
	"net/http"
	"os"
	"strconv"
	"time"
)

//TODO: исправить нейминг Wall -> Post

type WallHandler struct {
	service *service.WallService
}

func NewWallHandler(service *service.WallService) *WallHandler {
	return &WallHandler{
		service: service,
	}
}

func (p *WallHandler) Wall(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.Getpost(w, r)
	case http.MethodPost:
		p.CreatePost(w, r)
	case http.MethodDelete:
		p.DeletePost(w, r)
	case http.MethodPut:
		p.EditPost(w, r)
	default:
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
	}
}

func (p *WallHandler) Getpost(w http.ResponseWriter, r *http.Request) {

	UserID, _ := auth.GetUserId(r)

	post, err := p.service.GetPost(UserID)

	if err != nil {
		http.Error(w, "Ошибка при получении поста", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		log.Printf("Ошибка при получении данных %v", err)
	}

}

func (p *WallHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

	UserID, _ := auth.GetUserId(r)
	r.ParseMultipartForm(10 << 20)

	title := r.FormValue("title")
	text := r.FormValue("text")
	var img_scr = ""

	file, handler, err := r.FormFile("img")
	if err == nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)

		path := "./web/static/uploads/posts/" + filename

		dst, err := os.Create(path)
		if err != nil {
			http.Error(w, "Ошибка сохранения", 500)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println(err)
			return
		}

		img_scr = "/static/uploads/posts/" + filename
	}

	res := &domain.CreateWallRequest{
		Title: title,
		Text:  text,
		Img:   img_scr,
	}
	post, err := p.service.CreatePost(res, UserID)

	if err != nil {
		http.Error(w, "Ошибка при создании поста", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(post)

}

func (p *WallHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("id")

	PostID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	UserID, _ := auth.GetUserId(r)

	p.service.DeletePost(PostID, UserID)
}

func (p *WallHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	var Post *domain.CreateWallRequest

	err := json.NewDecoder(r.Body).Decode(&Post)
	idSTR := r.URL.Query().Get("id")

	PostID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	err = p.service.EditPostByID(Post, PostID)
	if err != nil {
		log.Println("Ошибка при вызове сервиса EditPost")
	}
}
