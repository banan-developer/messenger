package transport

import (
	"encoding/json"
	"messenger_v2/internal/domain"
	"messenger_v2/internal/repository"
	"messenger_v2/pkg/auth"
	"net/http"
	"strconv"
)

type GroupHandler struct{ repo *repository.ChatRepo }

func NewGroupHandler(repo *repository.ChatRepo) *GroupHandler { return &GroupHandler{repo: repo} }

// Groups обрабатывает создание, просмотр, переименование и изменение состава групповых чатов.
//
// # API Контракт
//
//  Маршрут:     /api/groups
//  Авторизация: Требуется (сессионная кука)
//
//  GET    /api/groups                  — получить список групп текущего пользователя
//  GET    /api/groups?chat_id={id}     — получить участников группы
//  POST   /api/groups                  — создать новую группу
//  PATCH  /api/groups?chat_id={id}     — переименовать группу
//  PUT    /api/groups?chat_id={id}&user_id={id} — добавить участника в группу
//  DELETE /api/groups?chat_id={id}     — удалить группу
//  DELETE /api/groups?chat_id={id}&user_id={id} — удалить участника из группы
//
// # Параметры запроса
//
//  chat_id (int, необяз.) — id группы/чата
//  user_id (int, необяз.) — id пользователя для добавления/удаления в группу
//
// # Формат данных запроса
//
//  Content-Type: application/json
//  Тело для POST: JSON-объект с полями title и user_ids
func (h *GroupHandler) Groups(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch { h.rename(w,r); return }
	if r.Method == http.MethodGet {
		if r.URL.Query().Get("chat_id") == "" {
			groups, err := h.repo.GetGroupsForUser(func() int { id, _ := auth.GetUserId(r); return id }())
			if err != nil {
				http.Error(w, "Failed to get groups", 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(groups)
			return
		}
		h.members(w, r)
		return
	}
	if r.Method == http.MethodDelete && r.URL.Query().Get("user_id") == "" { h.deleteGroup(w,r); return }
	if r.Method == http.MethodPut || r.Method == http.MethodDelete {
		h.changeMember(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, _ := auth.GetUserId(r)
	var req domain.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		http.Error(w, "Invalid group data", http.StatusBadRequest)
		return
	}
	seen := map[int]bool{userID: true}
	ids := []int{userID}
	for _, id := range req.UserIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	if len(ids) < 2 || len(ids) > 50 {
		http.Error(w, "Group must contain from 2 to 50 members", http.StatusBadRequest)
		return
	}
	chatID, err := h.repo.CreateGroupChat(req.Title, "", ids)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"chat_id": chatID})
}

// deleteGroup удаляет группу чата, если текущий пользователь является её создателем.
//
// # API Контракт
//
//  Метод:       DELETE
//  Маршрут:     /api/groups
//  Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//  chat_id (int, обяз.) — id удаляемой группы
func (h *GroupHandler) deleteGroup(w http.ResponseWriter,r *http.Request){ chatID,_:=strconv.Atoi(r.URL.Query().Get("chat_id")); caller,_:=auth.GetUserId(r); ok,_:=h.repo.IsCreator(chatID,caller); if !ok {http.Error(w,"Forbidden",403);return}; if err:=h.repo.DeleteGroup(chatID);err!=nil{http.Error(w,"Failed",500);return}; w.WriteHeader(204) }

// rename переименовывает группу, если текущий пользователь является её создателем.
//
// # API Контракт
//
//  Метод:       PATCH
//  Маршрут:     /api/groups
//  Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//  chat_id (int, обяз.) — id группы
//
// # Формат данных запроса
//
//  Content-Type: application/json
//  Тело: JSON-объект с полем title
func (h *GroupHandler) rename(w http.ResponseWriter,r *http.Request){ chatID,_:=strconv.Atoi(r.URL.Query().Get("chat_id")); caller,_:=auth.GetUserId(r); ok,_:=h.repo.IsCreator(chatID,caller); if !ok {http.Error(w,"Forbidden",403);return}; var body struct{Title string `json:"title"`}; if json.NewDecoder(r.Body).Decode(&body)!=nil||body.Title==""{http.Error(w,"Invalid title",400);return}; if err:=h.repo.RenameGroup(chatID,body.Title);err!=nil{http.Error(w,"Failed",500);return}; w.WriteHeader(204) }

// changeMember добавляет или удаляет участника из группы.
//
// # API Контракт
//
//  Метод:       PUT или DELETE
//  Маршрут:     /api/groups
//  Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//  chat_id (int, обяз.) — id группы
//  user_id (int, обяз.) — id пользователя, которого добавляют/удаляют
//
// # Ответы
//
//  204 No Content — действие успешно выполнено
//  403 Forbidden  — текущий пользователь не имеет прав или пользователь не является другом
func (h *GroupHandler) changeMember(w http.ResponseWriter, r *http.Request) {
	chatID, e1 := strconv.Atoi(r.URL.Query().Get("chat_id"))
	userID, e2 := strconv.Atoi(r.URL.Query().Get("user_id"))
	if e1 != nil || e2 != nil {
		http.Error(w, "Invalid data", 400)
		return
	}
	callerID, _ := auth.GetUserId(r)
	allowed, err := h.repo.IsCreator(chatID, callerID)
	if err != nil || !allowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodPut {
		friends, err := h.repo.AreFriends(callerID, userID)
		if err != nil || !friends {
			http.Error(w, "Only friends can be added", http.StatusForbidden)
			return
		}
	}
	if r.Method == http.MethodPut {
		err = h.repo.AddMemberToGroup(chatID, userID)
	} else {
		err = h.repo.RemoveMemberFromGroup(chatID, userID)
	}
	if err != nil {
		http.Error(w, "Failed to change members", 400)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// members возвращает список участников группы.
//
// # API Контракт
//
//  Метод:       GET
//  Маршрут:     /api/groups
//  Авторизация: Требуется (сессионная кука)
//
// # Параметры запроса
//
//  chat_id (int, обяз.) — id группы, участников которой нужно получить
//
// # Формат ответа
//
//  Content-Type: application/json
//  Тело: JSON-массив участников группы
func (h *GroupHandler) members(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.Atoi(r.URL.Query().Get("chat_id"))
	if err != nil {
		http.Error(w, "Invalid chat", 400)
		return
	}
	userID, _ := auth.GetUserId(r)
	allowed, err := h.repo.IsParticipant(chatID, userID)
	if err != nil || !allowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	members, err := h.repo.GetGroupMembers(chatID)
	if err != nil {
		http.Error(w, "Failed to get members", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
