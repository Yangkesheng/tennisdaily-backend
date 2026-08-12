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

type UserSettingsHandler struct {
	userSettingsService *service.UserSettingsService
}

func NewUserSettingsHandler(userSettingsService *service.UserSettingsService) *UserSettingsHandler {
	return &UserSettingsHandler{userSettingsService: userSettingsService}
}

func (h *UserSettingsHandler) Get(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}

	settings, err := h.userSettingsService.Get(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, settings)
}

func (h *UserSettingsHandler) Update(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req model.UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}

	logger.Debug("PUT /api/user-settings start userID=%d", userID)
	settings, err := h.userSettingsService.Update(userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, 404, response.CodeNotFound, "not found")
			return
		}
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("PUT /api/user-settings success userID=%d", userID)
	response.OK(c, settings)
}
