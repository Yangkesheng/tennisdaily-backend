package handler

import (
	"github.com/gin-gonic/gin"
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"
)

type TennisStringHandler struct{ service *service.TennisStringService }

func NewTennisStringHandler(s *service.TennisStringService) *TennisStringHandler {
	return &TennisStringHandler{service: s}
}
func (h *TennisStringHandler) List(c *gin.Context) {
	var q model.StringLibraryQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	result, err := h.service.List(q)
	if err != nil {
		logger.Error("string library list: %v", err)
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, result)
}
func (h *TennisStringHandler) Brands(c *gin.Context) {
	result, err := h.service.Brands()
	if err != nil {
		logger.Error("string brands list: %v", err)
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, result)
}
