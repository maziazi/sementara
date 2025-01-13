package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"project_sprint/pkg/config"
	"strings"
	"time"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := config.LoadEnv().JWTSecret
	if secret == "" {
		fmt.Println("⚠️  WARNING: JWT_SECRET tidak terbaca, gunakan default untuk debugging!")
		secret = "default-secret-key"
	}
	return secret
}

func GenerateToken(email string, userId uint) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userID": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 1 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWTAuthMiddleware memvalidasi token JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{} // Tidak pakai pointer (&) agar bisa diisi
		fmt.Println("🔹 Token diterima:", tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("❌ Invalid Token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		fmt.Println("✅ Parsed claims:", claims)
		fmt.Println("aa", claims["email"].(string))
		// Ambil email dari token dan simpan di context
		if userID, ok := claims["userID"].(float64); ok {
			c.Set("userID", int(userID)) // Konversi ke int sebelum disimpan
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		c.Next()
	}
}
