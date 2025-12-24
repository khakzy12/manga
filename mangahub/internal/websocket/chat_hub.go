package socket

import (
	"github.com/gorilla/websocket"
)

// Client represents a single chat participant
type Client struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
}

// ChatMessage represents the JSON structure for messages
type ChatMessage struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan ChatMessage
	Register   chan *Client
	Unregister chan *Client
}

func NewChatHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan ChatMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
		}
	}
}
