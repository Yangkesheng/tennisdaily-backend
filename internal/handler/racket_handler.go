package handler

import (
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RacketHandler struct {
	racketService *service.RacketService
}

func NewRacketHandler(racketService *service.RacketService) *RacketHandler {
	return &RacketHandler{racketService: racketService}
}

func (h *RacketHandler) Brands(c *gin.Context) {
	logger.Debug("GET /api/racket-brands start")
	brands, err := h.racketService.Brands()
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/racket-brands success count=%d", len(brands))
	response.OK(c, brands)
}

func (h *RacketHandler) Series(c *gin.Context) {
	var query model.RacketSeriesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	logger.Debug("GET /api/racket-series start brandID=%d", query.BrandID)
	series, err := h.racketService.Series(query)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/racket-series success brandID=%d count=%d", query.BrandID, len(series))
	response.OK(c, series)
}

func (h *RacketHandler) Library(c *gin.Context) {
	var query model.RacketLibraryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	logger.Debug("GET /api/racket-library start brandID=%d seriesID=%d", query.BrandID, query.SeriesID)
	groups, err := h.racketService.Library(query)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/racket-library success brandID=%d seriesID=%d brandCount=%d", query.BrandID, query.SeriesID, len(groups))
	response.OK(c, groups)
}

func (h *RacketHandler) MyRackets(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/my-rackets start userID=%d", userID)
	rackets, err := h.racketService.MyRacketsForSession(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/my-rackets success userID=%d count=%d", userID, len(rackets))
	response.OK(c, rackets)
}

func (h *RacketHandler) Stats(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/rackets/stats start userID=%d", userID)

	stats, err := h.racketService.Stats(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/rackets/stats success userID=%d racketCount=%d", userID, stats.RacketCount)
	response.OK(c, stats)
}

func (h *RacketHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	includeRetired := c.Query("includeRetired") == "true"
	logger.Debug("GET /api/rackets start userID=%d includeRetired=%t", userID, includeRetired)

	rackets, err := h.racketService.List(userID, includeRetired)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/rackets success userID=%d count=%d", userID, len(rackets))
	response.OK(c, rackets)
}

func (h *RacketHandler) Selectable(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/rackets/selectable start userID=%d", userID)

	rackets, err := h.racketService.Selectable(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/rackets/selectable success userID=%d count=%d", userID, len(rackets))
	response.OK(c, rackets)
}

func (h *RacketHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/rackets start userID=%d", userID)

	var req model.CreateRacketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	racket, err := h.racketService.Create(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/rackets success userID=%d racketID=%d", userID, racket.ID)
	response.OK(c, racket)
}

func (h *RacketHandler) Detail(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/rackets/:id start userID=%d racketID=%d", userID, id)

	detail, err := h.racketService.Detail(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("GET /api/rackets/:id success userID=%d racketID=%d", userID, id)
	response.OK(c, detail)
}

func (h *RacketHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("PUT /api/rackets/:id start userID=%d racketID=%d", userID, id)

	var req model.UpdateRacketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	racket, err := h.racketService.Update(userID, id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("PUT /api/rackets/:id success userID=%d racketID=%d", userID, racket.ID)
	response.OK(c, racket)
}

func (h *RacketHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("DELETE /api/rackets/:id start userID=%d racketID=%d", userID, id)

	if err := h.racketService.Delete(userID, id); err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("DELETE /api/rackets/:id success userID=%d racketID=%d", userID, id)
	response.OK(c, gin.H{"deleted": true})
}

func (h *RacketHandler) SetPrimary(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/rackets/:id/set-primary start userID=%d racketID=%d", userID, id)

	racket, err := h.racketService.SetPrimary(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/rackets/:id/set-primary success userID=%d racketID=%d", userID, racket.ID)
	response.OK(c, racket)
}

func (h *RacketHandler) Retire(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/rackets/:id/retire start userID=%d racketID=%d", userID, id)

	racket, err := h.racketService.Retire(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/rackets/:id/retire success userID=%d racketID=%d", userID, racket.ID)
	response.OK(c, racket)
}

func (h *RacketHandler) CreateStringingRecord(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	racketID, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/rackets/:id/stringing-records start userID=%d racketID=%d", userID, racketID)

	var req model.CreateStringingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	record, err := h.racketService.CreateStringingRecord(userID, racketID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/rackets/:id/stringing-records success userID=%d racketID=%d recordID=%d", userID, racketID, record.ID)
	response.OK(c, record)
}
