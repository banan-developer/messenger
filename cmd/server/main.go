package main

import (
	"database/sql"
	"fmt"
	"log"
	"messenger_v2/internal/repository"
	"messenger_v2/internal/service"
	"messenger_v2/internal/transport"
	"messenger_v2/internal/transport/websocket"
	"messenger_v2/pkg/auth"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// создание файлa, для отлатки ошибок и инфо
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// переменное окружение на пароля от бд
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		errorLog.Fatal("DB_PASSWORD не установлен")
	}

	// строка подключения к бд
	// dsn := fmt.Sprintf("root:%s@tcp(db:3306)/messanger", dbPassword)
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/messanger", dbPassword)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errorLog.Fatal("ошибка подлкючения к базе данных", err)
	}

	err = db.Ping()
	if err != nil {
		errorLog.Fatal("Ошибка подключения", err)
	}

	defer db.Close()
	if err := ensureAttachmentColumns(db); err != nil {
		errorLog.Fatal("Ошибка обновления схемы вложений: ", err)
	}

	auth.InitStore()
	hub := websocket.NewHub()

	// Инициализация слоев
	UserRepo := repository.NewUserRepository(db)
	UserService := service.NewUserService(UserRepo)
	UserHanlder := transport.NewUserHandler(UserService)

	AuthService := service.NewAuthService(UserRepo)
	AuthHandler := transport.NewAuthHandler(AuthService)

	WallRepo := repository.NewWallRepo(db)
	WallService := service.NewWallService(WallRepo)
	WallHanlder := transport.NewWallHandler(WallService)

	FriendRepo := repository.NewFriendRepo(db)
	FriendService := service.NewFrinedService(FriendRepo)
	FriendHandler := transport.NewFriendHandler(FriendService)

	MessagesRepo := repository.NewMessageRepo(db)
	ChatRepo := repository.NewChatRepo(db)
	MessagesService := service.NewMessageService(MessagesRepo)
	MessgesHandler := transport.NewMessageHandler(MessagesService, hub)
	GroupHandler := transport.NewGroupHandler(ChatRepo)

	fileServer := http.FileServer(http.Dir("./web/static"))

	http.Handle("/static/",
		http.StripPrefix("/static/", fileServer))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))

	http.Handle("/api/profile", auth.RequireAuth(http.HandlerFunc(UserHanlder.Profile)))
	http.Handle("/api/profile/avatar", auth.RequireAuth(http.HandlerFunc(UserHanlder.UploadAvatarUser)))
	http.Handle("/api/post", auth.RequireAuth(http.HandlerFunc(WallHanlder.Post)))
	http.Handle("/api/friend", auth.RequireAuth(http.HandlerFunc(FriendHandler.Friends)))
	http.Handle("/api/incomingrequest", auth.RequireAuth(http.HandlerFunc(FriendHandler.GetIncomigRequest)))
	http.Handle("/api/outgoingrequest", auth.RequireAuth(http.HandlerFunc(FriendHandler.OutgoingRequests)))
	http.Handle("/api/messages", auth.RequireAuth(http.HandlerFunc(MessgesHandler.Messages)))
	http.Handle("/api/messages/file", auth.RequireAuth(http.HandlerFunc(MessgesHandler.SendMessageWithFile)))
	http.Handle("/api/groups", auth.RequireAuth(http.HandlerFunc(GroupHandler.Groups)))

	wsHandler := transport.NewWebSocketHandler(MessagesService, hub)

	// Роут для WebSocket
	http.HandleFunc("/ws", wsHandler.HandleWS)

	http.HandleFunc("/login", AuthHandler.Login)
	http.HandleFunc("/registration", AuthHandler.Registration)
	http.HandleFunc("/exit", AuthHandler.Logout)

	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/friend", friendHandler)
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/friends", friendList)
	http.HandleFunc("/messages", messagesList)

	fmt.Println("Сервер запущен на http://127.0.0.1:8020/login")
	http.ListenAndServe(":8020", nil)
}

func ensureAttachmentColumns(db *sql.DB) error {
	columns := map[string]string{
		"attachment_name": "VARCHAR(255) NULL",
		"attachment_type": "VARCHAR(120) NULL",
		"attachment_size": "BIGINT NOT NULL DEFAULT 0",
	}

	for name, definition := range columns {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*)
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = 'messeges'
			  AND COLUMN_NAME = ?
		`, name).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			if _, err := db.Exec("ALTER TABLE messeges ADD COLUMN " + name + " " + definition); err != nil {
				return err
			}
		}
	}

	return nil
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/profile.html")
}

func friendHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/friend.html")
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/chat.html")
}

func friendList(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/friends.html")
}

func messagesList(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/messages.html")
}
