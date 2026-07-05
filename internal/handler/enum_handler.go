package handler

import (
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type EnumHandler struct {
	enumService *service.EnumService
}

func NewEnumHandler(enumService *service.EnumService) *EnumHandler {
	return &EnumHandler{enumService: enumService}
}

func (h *EnumHandler) SessionConfig(c *gin.Context) {
	logger.Debug("GET /api/session-config start")
	config := h.enumService.SessionConfig()
	logger.Debug("GET /api/session-config success")
	response.OK(c, config)
}
