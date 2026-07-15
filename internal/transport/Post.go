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

type PostHandler struct {
	service *service.PostService
}

func NewWallHandler(service *service.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

// Post является маршрутизатором для работы со стеной пользователя.
//
// # API Контракт
//
//	Маршрут:     /api/post
//	Авторизация: Требуется (сессионная кука)
//
//	GET    /api/post              — получить посты пользователя
//	GET    /api/post?user_id={id} — получить посты другого пользователя
//	POST   /api/post              — создать новый пост
//	PUT    /api/post?id={id}      — отредактировать пост
//	DELETE /api/post?id={id}      — удалить пост
func (p *PostHandler) Post(w http.ResponseWriter, r *http.Request) {
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

// Getpost возвращает список постов указанного пользователя или текущего пользователя из сессии.
//
// # API Контракт
//
//	Метод:       GET
//	Маршрут:     /api/post
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	user_id (int, необяз.) — id владельца стены. Если не передан, используются посты текущего пользователя.
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON-массив постов
func (p *PostHandler) Getpost(w http.ResponseWriter, r *http.Request) {

	var UserID int
	var err error
	UserIDSTR := r.URL.Query().Get("user_id")
	if UserIDSTR != "" {
		UserID, err = strconv.Atoi(UserIDSTR)
		if err != nil {
			http.Error(w, "Invalid note id", http.StatusBadRequest)
			return
		}

	} else {
		UserID, _ = auth.GetUserId(r)
	}

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

// CreatePost создает новый пост на стене текущего пользователя с возможностью загрузки изображения.
//
// # API Контракт
//
//	Метод:       POST
//	Маршрут:     /api/post
//	Авторизация: Требуется (сессионная кука)
//
// # Формат данных запроса
//
//	Формат: multipart/form-data
//	title (string, обяз.) — заголовок поста
//	text  (string, обяз.) — текст поста
//	img   (file, необяз.) — изображение поста
//
// # Формат ответа
//
//	Content-Type: application/json
//	Тело: JSON-массив обновленных постов пользователя
func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

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

// DeletePost удаляет пост по его идентификатору.
//
// # API Контракт
//
//	Метод:       DELETE
//	Маршрут:     /api/post
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	id (int, обяз.) — id поста, который нужно удалить
func (p *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idSTR := r.URL.Query().Get("id")

	PostID, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "Invalid note id", http.StatusBadRequest)
		return
	}
	UserID, _ := auth.GetUserId(r)

	p.service.DeletePost(PostID, UserID)
}

// EditPost обновляет существующий пост текущего пользователя.
//
// # API Контракт
//
//	Метод:       PUT
//	Маршрут:     /api/post
//	Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//	id (int, обяз.) — id поста, который нужно отредактировать
//
// # Формат данных запроса
//
//	Content-Type: application/json
//	Тело: JSON-объект с полями title и text
func (p *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
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
