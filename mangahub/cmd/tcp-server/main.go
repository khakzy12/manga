package main

import (
	"log"
	"mangahub/internal/tcp" // Ensure this matches your new structure
)

func main() {
	// 1. Initialize the TCP Server logic
	// We pass the port we want it to listen on
	server := tcp.NewProgressSyncServer("8081")

	log.Println("🛰️ Starting Standalone TCP Progress Sync Server...")

	// 2. Start the server (This is a blocking call, so no 'go' keyword here)
	server.Start()
}
