package handler

import (
	"errors"
	"strconv"

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

	sessions, err := h.sessionService.List(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, sessions)
}

func (h *SessionHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

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

	session, err := h.sessionService.FindByID(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
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

	if err := h.sessionService.Delete(userID, id); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (h *SessionHandler) Latest(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	session, err := h.sessionService.Latest(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	response.OK(c, session)
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
