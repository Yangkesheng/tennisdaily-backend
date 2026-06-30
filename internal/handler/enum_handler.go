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

func (h *EnumHandler) All(c *gin.Context) {
	logger.Debug("GET /api/enums start")
	enums := h.enumService.All()
	logger.Debug("GET /api/enums success")
	response.OK(c, enums)
}
