package handler

import (
	"errors"

	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/middleware"
	"tennisdaily-backend/internal/model"
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
	logger.Debug("GET /api/stats/month start userID=%d", userID)

	stats, err := h.statsService.Month(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/stats/month success userID=%d monthCount=%d totalCount=%d", userID, stats.MonthCount, stats.TotalCount)
	response.OK(c, stats)
}

func (h *StatsHandler) Charts(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}

	var query model.StatsChartsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}

	logger.Debug("GET /api/stats/charts start userID=%d period=%s year=%d month=%d", userID, query.Period, query.Year, query.Month)
	charts, err := h.statsService.Charts(userID, query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
			return
		}
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/stats/charts success userID=%d period=%s year=%d month=%d sessionCount=%d totalCost=%.2f", userID, charts.Period, charts.Year, charts.Month, charts.Summary.SessionCount, charts.Summary.TotalCost)
	response.OK(c, charts)
}

func (h *StatsHandler) Records(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}
	logger.Debug("GET /api/stats/records start userID=%d", userID)

	records, err := h.statsService.Records(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/stats/records success userID=%d totalCount=%d totalMinutes=%d streakDays=%d", userID, records.TotalCount, records.TotalMinutes, records.CurrentStreakDays)
	response.OK(c, records)
}
