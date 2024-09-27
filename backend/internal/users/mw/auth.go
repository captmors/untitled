package mw

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMWConfig struct {
	JwtKey []byte
	Claims func() jwt.Claims
}

var authUrls = []string{
	"/auth/",
}

func AuthMW(config AuthMWConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, url := range authUrls {
			if strings.HasPrefix(c.Request.URL.Path, url) {
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := config.Claims()

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := claims.(*jwt.RegisteredClaims); ok {
			c.Set("UserID", claims.Subject)
		}

		c.Next()
	}
}

func MockAuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("email", "mock@mock.com")
		c.Next()
	}
}
