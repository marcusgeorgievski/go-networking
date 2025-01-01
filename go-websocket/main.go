package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()

	mux := http.NewServeMux()

	mux.Handle("/ws", websocket.Handler(server.handleWebSocket))

	fmt.Println("Listening on ws://localhost:3000")

	http.ListenAndServe(":3000", mux)
}
