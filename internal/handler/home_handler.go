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

type HomeHandler struct {
	homeService *service.HomeService
}

func NewHomeHandler(homeService *service.HomeService) *HomeHandler {
	return &HomeHandler{homeService: homeService}
}

func (h *HomeHandler) Summary(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}

	var query model.HomeSummaryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}

	logger.Debug("GET /api/home/summary start userID=%d year=%d month=%d", userID, query.Year, query.Month)

	summary, err := h.homeService.Summary(userID, query)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
			return
		}
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}

	logger.Debug("GET /api/home/summary success userID=%d year=%d month=%d monthCount=%d totalCost=%.2f", userID, summary.Year, summary.Month, summary.Session.MonthCount, summary.Expense.TotalCost)
	response.OK(c, summary)
}
