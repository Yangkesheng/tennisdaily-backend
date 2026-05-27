package handler

import (
	"tennisdaily-backend/internal/middleware"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) Month(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}

	stats, err := h.statsService.Month(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, stats)
}
