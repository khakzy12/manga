package auth

import (
	"fmt"
	"net/http"

	//"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		source := "header"
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// Fallback: allow token via query param (useful for browser WebSocket handshakes)
			tokenString = c.Query("token")
			source = "query"
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "No token provided"})
			return
		}

		if source == "header" {
			// If the user sends "Bearer " with nothing after it
			if !strings.HasPrefix(tokenString, "Bearer ") || len(tokenString) < 8 {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Authorization header format"})
				return
			}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JWTKey, nil // Use the shared variable
		})

		if err != nil || !token.Valid {
			fmt.Printf("❌ JWT Error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Store user data in context for other handlers to use
			c.Set("user_id", fmt.Sprintf("%v", claims["user_id"]))
			c.Set("username", fmt.Sprintf("%v", claims["username"]))
			c.Set("role", fmt.Sprintf("%v", claims["role"]))
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		}
	}
}
