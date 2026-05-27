package middleware

import (
	"strings"

	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

const CurrentUserIDKey = "currentUserID"

func Auth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		userID, err := authService.ParseToken(token)
		if err != nil {
			response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set(CurrentUserIDKey, userID)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(CurrentUserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := value.(int64)
	return userID, ok
}
