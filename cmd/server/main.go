package main

import (
	"database/sql"
	"fmt"
	"log"
	"messenger_v2/internal/repository"
	"messenger_v2/internal/service"
	"messenger_v2/internal/transport"
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

	auth.InitStore()

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
	MessagesService := service.NewMessagesService(MessagesRepo)
	MessgesHandler := transport.NewMessageHandler(MessagesService)

	fileServer := http.FileServer(http.Dir("./web/static"))

	http.Handle("/static/",
		http.StripPrefix("/static/", fileServer))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))

	http.Handle("/api/profile", auth.RequireAuth(http.HandlerFunc(UserHanlder.Profile)))
	http.Handle("/api/profile/avatar", auth.RequireAuth(http.HandlerFunc(UserHanlder.UploadAvatarUser)))
	http.Handle("/api/post", auth.RequireAuth(http.HandlerFunc(WallHanlder.Post)))
	http.Handle("/api/friend", auth.RequireAuth(http.HandlerFunc(FriendHandler.Friends)))
	http.Handle("/api/messages", auth.RequireAuth(http.HandlerFunc(MessgesHandler.Messages)))

	http.HandleFunc("/login", AuthHandler.Login)
	http.HandleFunc("/registration", AuthHandler.Registration)
	http.HandleFunc("/exit", AuthHandler.Logout)

	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/friend", friendHandler)
	http.HandleFunc("/chat", chatHandler)

	fmt.Println("Сервер запущен на http://127.0.0.1:8020/login")
	http.ListenAndServe(":8020", nil)
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
