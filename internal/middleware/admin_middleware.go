package middleware

import (
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// RequireAdmin 校验当前登录用户是否在管理员白名单内，非管理员返回 403。
func RequireAdmin(adminUserIDs []int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := CurrentUserID(c)
		if !ok || !service.IsAdminUser(userID, adminUserIDs) {
			response.Error(c, 403, response.CodeForbidden, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
