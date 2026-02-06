package main

import (
	"go-study-case3/internal/ws"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &ws.Client{
		Hub: hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	
	client.WritePump()
	client.ReadPump()
}

func main() {
	hub := ws.NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", nil)
	err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}