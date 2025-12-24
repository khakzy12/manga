package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mangahub/pkg/models"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	DB            *sql.DB
	TCPServerAddr string
}

// POST /users/library
func (uc *UserController) AddToLibrary(c *gin.Context) {
	// Retrieve UserID from JWT (set in middleware)
	userID, _ := c.Get("user_id")

	var input struct {
		MangaID string `json:"manga_id" binding:"required"`
		Status  string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := uc.DB.Exec(`
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status) 
		VALUES (?, ?, 0, ?)
		ON CONFLICT(user_id, manga_id) DO UPDATE SET status=excluded.status`,
		userID, input.MangaID, input.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update library"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Library updated"})
}

// UpdateProgress handles the PUT request and triggers the TCP broadcast
func (uc *UserController) UpdateProgress(c *gin.Context) {
	var input models.ProgressUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Ensure we are assigning to Chapter
	uname, _ := c.Get("username")
	input.Username = fmt.Sprintf("%v", uname)

	// Dial TCP and send...
	conn, err := net.Dial("tcp", uc.TCPServerAddr)
	if err == nil {
		defer conn.Close()
		json.NewEncoder(conn).Encode(input)
	}

	c.JSON(200, gin.H{"message": "Chapter progress updated"})
}
