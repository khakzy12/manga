package auth

import (
	"database/sql"
	"fmt"
	"mangahub/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTKey = []byte("MangaHub_Secret_Key_2024")

type AuthController struct {
	DB *sql.DB
}

// Register Request Structure
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	// ... binding ...

	// 1. Hash the password before saving!
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 2. Insert using the CORRECT column name
	_, err := ac.DB.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		input.Username, string(hashed), "user")

	if err != nil {
		// Log the ACTUAL error to your terminal so you can see it
		fmt.Println("Registration Error:", err)
		c.JSON(400, gin.H{"error": "Registration failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Success"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input RegisterInput // Reusing same structure for login
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("❌ JSON Binding Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
