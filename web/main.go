package main

import (
	"database/sql"
	"fmt"
	"log"
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
	Name string `json:"name"`
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

	http.HandleFunc("/ws", app.wsHandler)
	http.HandleFunc("/test", app.getname)
	http.HandleFunc("/", app.HomeHandler)

	fileServer := http.FileServer(http.Dir("./pkg/ui/static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	go hanldeMessage()
	fmt.Println("Сервер запущен на http://127.0.0.1:4040")
	http.ListenAndServe(":4040", nil)
}
