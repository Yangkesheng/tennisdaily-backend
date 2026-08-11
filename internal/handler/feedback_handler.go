package handler

import (
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FeedbackHandler struct {
	feedbackService *service.FeedbackService
}

func NewFeedbackHandler(feedbackService *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{feedbackService: feedbackService}
}

func (h *FeedbackHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/feedback start userID=%d", userID)

	var req model.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	feedback, err := h.feedbackService.Create(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/feedback success userID=%d feedbackID=%d", userID, feedback.ID)
	response.OK(c, feedback)
}
