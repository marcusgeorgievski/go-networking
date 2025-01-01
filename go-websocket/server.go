package main

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	connections map[*websocket.Conn]bool
	mu sync.Mutex
}

const MAX_BYTES = 1024

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

// 
func (s *Server) handleWebSocket(ws *websocket.Conn) {
	s.mu.Lock()
	s.connections[ws] = true
	s.mu.Unlock()

	if len(s.connections) == 0 {
		fmt.Println("No connections.")
	}

	fmt.Println("Connection added:", ws.RemoteAddr())
	s.displayConnections()

	s.readLoop(ws)
}

// 
func (s *Server) removeConnection(ws *websocket.Conn) {
	s.mu.Lock()
	delete(s.connections, ws)
	s.mu.Unlock()

	if err := ws.Close(); err != nil {
		fmt.Println("Error closing connection:", err)
	}

	fmt.Println("Connection removed:", ws.RemoteAddr())
	s.displayConnections()
}

// 
func (s *Server) readLoop(ws *websocket.Conn) {
	defer s.removeConnection(ws)

	buf := make([]byte, MAX_BYTES)

	for {
		n, err := ws.Read(buf)
		
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client:", ws.RemoteAddr())
				break
			}
			fmt.Println("Error reading:", err)
			continue
		}

		msg := buf[:n]
		fmt.Println("Message:", string(msg))
		s.broadcast(msg, ws)
	}
}



func (s *Server) broadcast(msg []byte, exclude *websocket.Conn) {
	s.mu.Lock()
    defer s.mu.Unlock()

	for conn := range s.connections {
		if conn == exclude {
			continue
		}

		go func(c *websocket.Conn) {
			if _, err := c.Write(msg); err != nil {
				fmt.Println("Broadcast error:", err)
			}
		}(conn)
	}
}

func (s *Server) displayConnections() {
	s.mu.Lock()
    defer s.mu.Unlock()

	fmt.Println("Connections:")
	for conn := range s.connections {
		fmt.Println(" -", conn.RemoteAddr())
	}
}
