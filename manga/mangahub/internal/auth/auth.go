package auth

import (
	"fmt"
	//"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString = c.Query("token")
		} else {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "No token provided"})
			return
		}

		// PARSING LOGIC
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm is HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// RETURN THE EXACT SAME KEY USED IN LOGIN
			return []byte("MangaHub_Secret_Key_2024"), nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("❌ JWT Error: %v\n", err) // DEBUG: Check terminal for this
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// RESOLVING CLAIMS
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Save as strings to be safe for the WebSocket logic
			c.Set("user_id", fmt.Sprintf("%v", claims["user_id"]))
			c.Set("username", fmt.Sprintf("%v", claims["username"]))
			c.Set("role", fmt.Sprintf("%v", claims["role"]))
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid claims"})
		}
	}
}
