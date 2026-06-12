package main

import (
	"database/sql"
	"fmt"
	"log"
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

	fileServer := http.FileServer(http.Dir("./web/static/css"))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	http.Handle("/static/", http.StripPrefix("/static/css", fileServer))

	http.HandleFunc("/profile", profileHandler)

	fmt.Println("Сервер запущен на http://127.0.0.1:4562/login")
	http.ListenAndServe(":4562", nil)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/html/profile.html")
}
