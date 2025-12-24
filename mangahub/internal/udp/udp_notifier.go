package udp

import (
	"fmt"
	socket "mangahub/internal/websocket"
	"net"
)

type NotificationServer struct {
	Port string
	Hub  *socket.Hub // Add reference to the Chat Hub
}

func (s *NotificationServer) Start() {
	addr, err := net.ResolveUDPAddr("udp", ":"+s.Port)
	if err != nil {
		fmt.Println("UDP Resolve Error:", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("UDP Listen Error:", err)
		return
	}
	fmt.Println("📣 UDP Notification Server listening on port", s.Port)

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("UDP Read Error:", err)
			continue
		}

		// Convert UDP byte data to string message
		receivedMsg := string(buf[:n])
		fmt.Printf("☁️ UDP Received: %s\n", receivedMsg)

		// BROADCAST TO WEBSOCKET USERS
		// We send this to the Hub so it appears in the browser chat
		if s.Hub != nil {
			s.Hub.Broadcast <- socket.ChatMessage{
				Username: "SYSTEM-BROADCAST",
				Message:  receivedMsg,
			}
		}
	}
}
