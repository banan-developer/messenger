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
	"time"
)

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
		log.Println("Save path:", path)
		if err != nil {
			http.Error(w, "Ошибка сохранения", 500)
			return
		}
		defer dst.Close()

		n, err := io.Copy(dst, file)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Copied bytes:", n)

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
