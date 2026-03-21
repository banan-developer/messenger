package main

import (
	"database/sql"
	"fmt"
	"log"
	"messenger/auth"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *sql.DB
}

type person struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	About    string `json:"about"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
}

type wall struct {
	Id    int    `json:"idwall"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Img   string `json:"img"`
}

func main() {
	// создание файле, для отлатки ошибок и инфо
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		db:       db,
	}

	auth.InitStore()

	http.Handle("/ws", auth.RequireAuth(http.HandlerFunc(app.wsHandler)))
	http.Handle("/api/profile", auth.RequireAuth(http.HandlerFunc(app.handleProfile)))
	http.Handle("/api/post", auth.RequireAuth(http.HandlerFunc(app.handlePost)))
	http.Handle("/api/friend", auth.RequireAuth(http.HandlerFunc(app.handleFriend)))
	http.Handle("/api/profile/avatar", auth.RequireAuth(http.HandlerFunc(app.handleProfileAvatar)))

	http.HandleFunc("/login", app.autoresHandler)
	http.HandleFunc("/register", app.regHanlder)
	http.HandleFunc("/exit", app.exitSession)

	fileServer := http.FileServer(http.Dir("./pkg/ui/static"))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./pkg/ui/js"))))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	// http.Handle("/", auth.RequireAuth(http.HandlerFunc(app.HomeHandler)))
	http.HandleFunc("/", app.HomeHandler)
	http.HandleFunc("/profile", app.HomeHandler)
	http.Handle("/anotherProfile.html", auth.RequireAuth(http.HandlerFunc(app.anotherProfilePage)))
	http.Handle("/chat.html", auth.RequireAuth(http.HandlerFunc(app.chatWithFriend)))

	go hanldeMessage()
	fmt.Println("Сервер запущен на http://127.0.0.1:8080/login")
	http.ListenAndServe(":8080", nil)
}
