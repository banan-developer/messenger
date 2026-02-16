package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	fmt.Println("Запуск сервера на http://127.0.0.1:4040/")

	http.ListenAndServe(":4040", mux)
}
