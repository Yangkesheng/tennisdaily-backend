package handler

import (
	"errors"

	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) WechatLogin(c *gin.Context) {
	logger.Debug("POST /api/auth/wechat-login start clientIP=%s", c.ClientIP())

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}

	loginResp, err := h.authService.WechatLogin(req.Code)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
			return
		}
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}

	logger.Debug("POST /api/auth/wechat-login success userID=%d openid=%s", loginResp.UserID, loginResp.OpenID)
	response.OK(c, loginResp)
}
