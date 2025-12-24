package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"mangahub/pkg/models"
	"net"
)

type ProgressSyncServer struct {
	Port string
}

func NewProgressSyncServer(port string) *ProgressSyncServer {
	return &ProgressSyncServer{Port: port}
}

func (s *ProgressSyncServer) Start() {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatalf("❌ TCP Error: %v", err)
	}
	defer listener.Close()

	fmt.Printf("✅ TCP Server active on port %s\n", s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("⚠️ Connection error: %v", err)
			continue
		}
		// Handle each sync request in a new goroutine
		go s.handleSync(conn)
	}
}

func (s *ProgressSyncServer) handleSync(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var update models.ProgressUpdate // Using shared model

	if err := decoder.Decode(&update); err == nil {
		fmt.Printf("🔄 [TCP SYNC] User %s is on Chapter %s\n", update.Username, update.Chapter)
	}
}
