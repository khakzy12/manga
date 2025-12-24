package auth

import (
	"database/sql"
	"fmt"
	"mangahub/pkg/models"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var JWTKey = []byte("MangaHub_Secret_Key_2024")

// UDP notification address used to broadcast system messages (overrides allowed for tests)
var UDPNotifyAddr = "127.0.0.1:12345"

type AuthController struct {
	DB *sql.DB
}

// Register Request Structure
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("❌ Register JSON Bind Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "detail": err.Error()})
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	if len(input.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters"})
		return
	}
	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	// Check if username exists
	var existingID string
	err := ac.DB.QueryRow("SELECT id FROM users WHERE username = ?", input.Username).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	if err != sql.ErrNoRows && err != nil {
		fmt.Println("DB Error checking user exists:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Hash Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	id := uuid.NewString()
	_, err = ac.DB.Exec("INSERT INTO users (id, username, password_hash, role) VALUES (?, ?, ?, ?)",
		id, input.Username, string(hashed), "user")
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(strings.ToLower(err.Error()), "constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
		fmt.Println("Registration Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	// Broadcast a system message via UDP to notify other services (e.g., WebSocket hub)
	go func(username string) {
		addr, err := net.ResolveUDPAddr("udp", UDPNotifyAddr)
		if err != nil {
			fmt.Println("UDP Resolve Error:", err)
			return
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("UDP Dial Error:", err)
			return
		}
		defer conn.Close()
		payload := []byte("📢 NEW USER: " + username)
		if _, err := conn.Write(payload); err != nil {
			fmt.Println("UDP Send Error:", err)
		}
	}(input.Username)

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "id": id})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input RegisterInput // Reusing same structure for login
	if err := c.ShouldBindJSON(&input); err != nil {
		// Add this print to see the ACTUAL error in your terminal
		fmt.Println("Binding Error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := ac.DB.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username = ?",
		input.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)

	if err != nil {
		fmt.Println("❌ Database Query Error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  fmt.Sprint(user.ID), // Is this a string or int?
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(JWTKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}
