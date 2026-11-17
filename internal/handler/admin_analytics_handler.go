package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/middleware"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"
)

type AdminAnalyticsHandler struct {
	service       *service.AdminAnalyticsService
	installations *service.AnalyticsInstallationService
}

func NewAdminAnalyticsHandler(s *service.AdminAnalyticsService, installations *service.AnalyticsInstallationService) *AdminAnalyticsHandler {
	return &AdminAnalyticsHandler{service: s, installations: installations}
}

func (h *AdminAnalyticsHandler) AppOpen(c *gin.Context) {
	var req model.AppOpenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	if err := h.installations.Touch(req.InstallationID); err != nil {
		logger.Error("POST /api/analytics/app-open failed: %v", err)
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, gin.H{"accepted": true})
}

func (h *AdminAnalyticsHandler) IdentifyInstallation(c *gin.Context) {
	var req model.AppOpenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}
	if err := h.installations.Identify(req.InstallationID, userID); err != nil {
		logger.Error("POST /api/analytics/app-open/identify failed: %v", err)
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, gin.H{"identified": true})
}

func (h *AdminAnalyticsHandler) Daily(c *gin.Context) {
	var query struct {
		StartDate string `form:"startDate" binding:"required"`
		EndDate   string `form:"endDate" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	result, err := h.service.Daily(query.StartDate, query.EndDate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
			return
		}
		logger.Error("GET /api/admin/analytics/daily failed: %v", err)
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/admin/analytics/daily success days=%d", len(result.Days))
	response.OK(c, result)
}
