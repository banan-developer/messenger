package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var client = make(map[*websocket.Conn]bool) // список всех подключенных пользователей
var broadcast = make(chan []byte)           // канал сообщений

// функция для перехода с http на websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func hanldeMessage() {
	for {
		msg := <-broadcast
		for person := range client {
			person.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
