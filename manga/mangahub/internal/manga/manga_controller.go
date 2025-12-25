package manga

import (
	"context"
	"mangahub/proto"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MangaController struct {
	GRPCClient proto.MangaServiceClient
}

// GET /manga/:id
func (mc *MangaController) GetMangaDetails(c *gin.Context) {
	id := c.Param("id")
	println("🔍 API Gateway: Searching for manga with ID:", id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Calling the gRPC Internal Service (Requirement 5)
	println("📡 API Gateway: Calling gRPC server...")
	resp, err := mc.GRPCClient.GetManga(ctx, &proto.GetMangaRequest{Id: id})
	if err != nil {
		println("❌ API Gateway: gRPC error:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found", "details": err.Error()})
		return
	}

	println("✅ API Gateway: Found manga:", resp.Title)
	c.JSON(http.StatusOK, resp)
}

func (mc *MangaController) AddManga(c *gin.Context) {
	// 1. Verify Admin Role (from JWT middleware)
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admins only"})
		return
	}

	// 2. Add to DB/gRPC... (logic here)

	// 3. BROADCAST via UDP
	udpAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255:12345")
	conn, _ := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	msg := []byte("New Manga Added: One Piece Vol 100!")
	conn.Write(msg)
}
