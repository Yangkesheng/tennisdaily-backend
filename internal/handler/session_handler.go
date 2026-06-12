package handler

import (
	"errors"
	"strconv"

	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/middleware"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	date := c.Query("date")
	logger.Debug("GET /api/sessions start userID=%d date=%s", userID, date)

	var sessions []model.SessionResponse
	var err error
	if date == "" {
		sessions, err = h.sessionService.List(userID)
	} else {
		sessions, err = h.sessionService.ListByDate(userID, date)
	}
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("GET /api/sessions success userID=%d date=%s count=%d", userID, date, len(sessions))
	response.OK(c, sessions)
}

func (h *SessionHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/sessions start userID=%d", userID)

	var req model.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}

	session, err := h.sessionService.Create(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/sessions success userID=%d sessionID=%d", userID, session.ID)
	response.OK(c, session)
}

func (h *SessionHandler) Get(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/sessions/:id start userID=%d sessionID=%d", userID, id)

	session, err := h.sessionService.FindByID(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("GET /api/sessions/:id success userID=%d sessionID=%d", userID, session.ID)
	response.OK(c, session)
}

func (h *SessionHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("PUT /api/sessions/:id start userID=%d sessionID=%d", userID, id)

	var req model.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}

	session, err := h.sessionService.Update(userID, id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("PUT /api/sessions/:id success userID=%d sessionID=%d", userID, session.ID)
	response.OK(c, session)
}

func (h *SessionHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("DELETE /api/sessions/:id start userID=%d sessionID=%d", userID, id)

	if err := h.sessionService.Delete(userID, id); err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("DELETE /api/sessions/:id success userID=%d sessionID=%d", userID, id)
	response.OK(c, gin.H{"deleted": true})
}

func (h *SessionHandler) Latest(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/sessions/latest start userID=%d", userID)

	session, err := h.sessionService.Latest(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	if session == nil {
		logger.Debug("GET /api/sessions/latest success userID=%d sessionID=<nil>", userID)
	} else {
		logger.Debug("GET /api/sessions/latest success userID=%d sessionID=%d", userID, session.ID)
	}
	response.OK(c, session)
}

func (h *SessionHandler) Calendar(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	logger.Debug("GET /api/sessions/calendar start userID=%d year=%d month=%d", userID, year, month)

	calendar, err := h.sessionService.Calendar(userID, year, month)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("GET /api/sessions/calendar success userID=%d year=%d month=%d activeDayCount=%d", userID, year, month, calendar.ActiveDayCount)
	response.OK(c, calendar)
}

func currentUserID(c *gin.Context) (int64, bool) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
		return 0, false
	}
	return userID, true
}

func parseIDParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return 0, false
	}
	return id, true
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRequest):
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
	case errors.Is(err, service.ErrUnauthorized):
		response.Error(c, 401, response.CodeUnauthorized, "unauthorized")
	case errors.Is(err, service.ErrForbidden):
		response.Error(c, 403, response.CodeForbidden, "forbidden")
	case errors.Is(err, service.ErrNotFound):
		response.Error(c, 404, response.CodeNotFound, "not found")
	default:
		response.Error(c, 500, response.CodeInternalError, "internal error")
	}
}
