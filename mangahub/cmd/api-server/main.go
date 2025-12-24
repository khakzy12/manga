package main

import (
	"fmt"
	"log"
	"mangahub/internal/auth"
	"mangahub/internal/manga"
	"mangahub/internal/udp"
	"mangahub/internal/user"
	socket "mangahub/internal/websocket"
	"mangahub/pkg/database"
	"mangahub/proto"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func SendProgressToTCPServer(message string) {
	// Connect to the standalone TCP server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("❌ Could not connect to TCP Sync Server: %v", err)
		return
	}
	defer conn.Close()

	// Send the data
	conn.Write([]byte(message))
}

// Helper for Admin UDP Broadcast
func broadcastNewManga(message string) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	if err != nil {
		log.Printf("UDP Resolve Error: %v", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("UDP Dial Error: %v", err)
		return
	}
	defer conn.Close()

	payload := []byte("📢 ADMIN NOTIFICATION: " + message)
	_, err = conn.Write(payload)
	if err != nil {
		log.Printf("UDP Broadcast Error: %v", err)
	} else {
		log.Printf("🚀 UDP Broadcast sent: %s", message)
	}
}

func main() {
	// 1. Initialize Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("DB Error:", err)
	}

	// 2. Initialize gRPC Client
	gConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("gRPC Connection Error:", err)
	}
	defer gConn.Close()
	mangaClient := proto.NewMangaServiceClient(gConn)

	// 3. Start Background Servers (Hub, TCP, UDP)
	hub := socket.NewChatHub()
	go hub.Run()

	udpServer := udp.NotificationServer{
		Port: "12345",
		Hub:  hub, // Connect the hub to the UDP server here!
	}
	go udpServer.Start()

	// 4. Initialize Gin
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authCtrl := &auth.AuthController{DB: db}
	mangaCtrl := &manga.MangaController{GRPCClient: mangaClient}
	userCtrl := &user.UserController{
		DB:            db,
		TCPServerAddr: "localhost:8081", // Pointing to your standalone TCP server
	}

	// --- ROUTES ---

	// Public Routes
	r.POST("/auth/register", authCtrl.Register)
	r.POST("/auth/login", authCtrl.Login)

	protected := r.Group("/")
	protected.Use(auth.AuthRequired())
	{
		protected.POST("/users/library", userCtrl.AddToLibrary)
		protected.PUT("/users/progress", userCtrl.UpdateProgress)
	}
	r.GET("/manga/:id", mangaCtrl.GetMangaDetails)
	r.GET("/manga/search", mangaCtrl.SearchManga)

	r.GET("/debug/ids", func(c *gin.Context) {
		rows, _ := db.Query("SELECT id FROM manga")
		var ids []string
		for rows.Next() {
			var id string
			rows.Scan(&id)
			ids = append(ids, id)
		}
		c.JSON(200, ids)
	})

	// Admin Route (UDP Trigger)
	r.POST("/admin/add-manga", func(c *gin.Context) {
		// 1. Role Check
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(403, gin.H{"error": "Admin only"})
			return
		}

		var input struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Author string `json:"author"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// 2. Add to Database
		_, err := db.Exec("INSERT INTO manga (id, title, author) VALUES (?, ?, ?)",
			input.ID, input.Title, input.Author)
		if err != nil {
			c.JSON(500, gin.H{"error": "DB Error: " + err.Error()})
			return
		}

		// 3. BROADCAST: Notify the network via UDP
		// This will go to the UDP server, which sends it to the Hub
		broadcastNewManga("New Manga Added: " + input.Title)

		c.JSON(200, gin.H{"status": "Manga created and broadcast sent!"})
	})

	r.DELETE("/admin/manga/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM manga WHERE id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete"})
			return
		}
		c.JSON(200, gin.H{"message": "Manga removed"})
	})

	// WebSocket Route (REMOVED DUPLICATE - Keeping the Protected version)
	// This satisfies the "Distinguish UserID" requirement using JWT
	r.GET("/ws/guest", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &socket.Client{
			Conn:     conn,
			UserID:   "GUEST",
			Username: "Guest_Viewer",
		}
		hub.Register <- client

		// IMPORTANT: You need this loop to keep the connection alive!
		go func() {
			defer func() { hub.Unregister <- client }()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})

	r.GET("/ws/chat", auth.AuthRequired(), func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		uname, _ := c.Get("username")

		log.Printf("Connect Attempt: User %v", uname)

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WS Upgrade Error: %v", err)
			return
		}

		client := &socket.Client{
			Conn:     conn,
			UserID:   fmt.Sprintf("%v", uid),
			Username: fmt.Sprintf("%v", uname),
		}

		hub.Register <- client

		go func() {
			defer func() { hub.Unregister <- client }()
			for {
				var msg socket.ChatMessage

				if err := conn.ReadJSON(&msg); err != nil {
					log.Printf("Read Error: %v", err)
					break
				}

				msg.Username = client.Username
				msg.UserID = client.UserID
				hub.Broadcast <- msg
			}
		}()
	})

	log.Println("🚀 Gateway running on :8080")
	r.Run(":8080")
}
