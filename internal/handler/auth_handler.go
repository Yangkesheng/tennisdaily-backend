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
		logger.Error("POST /api/auth/wechat-login internal error: %v", err)
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}

	logger.Debug("POST /api/auth/wechat-login success userID=%d openid=%s", loginResp.UserID, loginResp.OpenID)
	response.OK(c, loginResp)
}

func (h *AuthHandler) PhoneLogin(c *gin.Context) {
	logger.Debug("POST /api/auth/phone-login start clientIP=%s", c.ClientIP())

	var req model.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "登录参数不完整")
		return
	}

	loginResp, err := h.authService.PhoneLogin(req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	logger.Debug("POST /api/auth/phone-login success userID=%d phone=%s isNewUser=%t", loginResp.User.ID, model.MaskPhone(loginResp.User.Phone), loginResp.IsNewUser)
	response.OK(c, loginResp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	user, err := h.authService.CurrentUser(userID)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.OK(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.OK(c, h.authService.Logout())
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.OK(c, user)
}

func handleAuthError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrInvalidRequest) {
		response.Error(c, 400, response.CodeInvalidRequest, "登录参数不完整")
		return
	}
	if errors.Is(err, service.ErrSensitiveContent) {
		response.Error(c, 400, response.CodeInvalidRequest, "输入内容包含敏感信息，请修改后重试")
		return
	}
	if errors.Is(err, service.ErrWechatLoginCode) {
		response.Error(c, 400, response.CodeInvalidRequest, "微信登录凭证无效，请重试")
		return
	}
	if errors.Is(err, service.ErrWechatPhoneCode) {
		response.Error(c, 400, response.CodeInvalidRequest, "手机号授权已失效，请重试")
		return
	}
	if errors.Is(err, service.ErrWechatService) {
		response.Error(c, 502, response.CodeInternalError, "微信服务暂时不可用，请稍后重试")
		return
	}
	if errors.Is(err, service.ErrUnauthorized) {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, 404, response.CodeNotFound, "not found")
		return
	}
	logger.Error("auth internal error path=%s error=%v", c.FullPath(), err)
	response.Error(c, 500, response.CodeInternalError, "internal error")
}
