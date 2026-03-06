package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./pkg/ui/html/profile.html")
}

// func (app *application) profileHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		app.getname(w, r)
// 	}

// }

func (app *application) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // переход из http в websocket
	if err != nil {
		fmt.Println("Ошибка при переходе на Websocket протокол")
		return
	}
	client[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("client disconnected")
			delete(client, conn)
			conn.Close()
			break
		}
		broadcast <- message

	}
}

func (app *application) getname(w http.ResponseWriter, r *http.Request) {
	var person person

	err := app.db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&person.Name)
	if err != nil {
		fmt.Println("error1!")
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(person)
	if err != nil {
		fmt.Println("error2!")
	}

}
