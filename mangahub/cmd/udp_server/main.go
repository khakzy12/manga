package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "127.0.0.1:12345"
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	// Default message or take from command line
	msg := "Admin added: Chainsaw Man Chapter 1"
	if len(os.Args) > 1 {
		msg = strings.Join(os.Args[1:], " ")
	}

	conn.Write([]byte(msg))
	fmt.Println("🚀 Admin Trigger: UDP Broadcast sent to network!")
}
