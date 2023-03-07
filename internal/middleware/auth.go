package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/tapiaw38/auth-api/internal/models"
	"github.com/tapiaw38/auth-api/internal/server"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
		"verify-email",
	}
)

// shouldCheckToken is a function that checks if the route should be checked for token
func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

// CheckAuthMiddleware is a middleware that checks if the user is authenticated
func CheckAuthMiddleware(s server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldCheckToken(c.Request.URL.Path) {
			c.Next()
			return
		}

		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}
