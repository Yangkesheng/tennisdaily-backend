package handler

import (
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ShoeHandler struct {
	shoeService *service.ShoeService
}

func NewShoeHandler(shoeService *service.ShoeService) *ShoeHandler {
	return &ShoeHandler{shoeService: shoeService}
}

func (h *ShoeHandler) Brands(c *gin.Context) {
	logger.Debug("GET /api/shoe-brands start")
	brands, err := h.shoeService.Brands()
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoe-brands success count=%d", len(brands))
	response.OK(c, brands)
}

func (h *ShoeHandler) Series(c *gin.Context) {
	var query model.ShoeSeriesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	logger.Debug("GET /api/shoe-series start brandID=%d gender=%v", query.BrandID, query.Gender)
	series, err := h.shoeService.Series(query)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoe-series success brandID=%d count=%d", query.BrandID, len(series))
	response.OK(c, series)
}

func (h *ShoeHandler) Library(c *gin.Context) {
	var query model.ShoeLibraryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	logger.Debug("GET /api/shoe-library start brandID=%d seriesID=%d gender=%v", query.BrandID, query.SeriesID, query.Gender)
	groups, err := h.shoeService.Library(query)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoe-library success brandID=%d seriesID=%d brandCount=%d", query.BrandID, query.SeriesID, len(groups))
	response.OK(c, groups)
}

func (h *ShoeHandler) LibraryStats(c *gin.Context) {
	logger.Debug("GET /api/shoe-library/stats start")
	stats, err := h.shoeService.LibraryStats()
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoe-library/stats success brandCount=%d", len(stats))
	response.OK(c, stats)
}

func (h *ShoeHandler) MyShoes(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/my-shoes start userID=%d", userID)
	shoes, err := h.shoeService.MyShoesForSession(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/my-shoes success userID=%d count=%d", userID, len(shoes))
	response.OK(c, shoes)
}

func (h *ShoeHandler) Stats(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/shoes/stats start userID=%d", userID)

	stats, err := h.shoeService.Stats(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoes/stats success userID=%d shoeCount=%d", userID, stats.ShoeCount)
	response.OK(c, stats)
}

func (h *ShoeHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	includeRetired := c.Query("includeRetired") == "true"
	logger.Debug("GET /api/shoes start userID=%d includeRetired=%t", userID, includeRetired)

	shoes, err := h.shoeService.List(userID, includeRetired)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoes success userID=%d count=%d", userID, len(shoes))
	response.OK(c, shoes)
}

func (h *ShoeHandler) Selectable(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/shoes/selectable start userID=%d", userID)

	shoes, err := h.shoeService.Selectable(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	logger.Debug("GET /api/shoes/selectable success userID=%d count=%d", userID, len(shoes))
	response.OK(c, shoes)
}

func (h *ShoeHandler) MyPrimaryShoe(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/my-shoes/primary start userID=%d", userID)

	shoe, err := h.shoeService.Primary(userID)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError, "internal error")
		return
	}
	if shoe == nil {
		logger.Debug("GET /api/my-shoes/primary success userID=%d shoeID=<nil>", userID)
	} else {
		logger.Debug("GET /api/my-shoes/primary success userID=%d shoeID=%d", userID, shoe.ID)
	}
	response.OK(c, shoe)
}

func (h *ShoeHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/shoes start userID=%d", userID)

	var req model.CreateShoeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	shoe, err := h.shoeService.Create(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/shoes success userID=%d shoeID=%d", userID, shoe.ID)
	response.OK(c, shoe)
}

func (h *ShoeHandler) Detail(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("GET /api/shoes/:id start userID=%d shoeID=%d", userID, id)

	shoe, err := h.shoeService.Detail(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("GET /api/shoes/:id success userID=%d shoeID=%d", userID, id)
	response.OK(c, shoe)
}

func (h *ShoeHandler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("PUT /api/shoes/:id start userID=%d shoeID=%d", userID, id)

	var req model.UpdateShoeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, response.CodeInvalidRequest, "invalid request")
		return
	}
	shoe, err := h.shoeService.Update(userID, id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("PUT /api/shoes/:id success userID=%d shoeID=%d", userID, shoe.ID)
	response.OK(c, shoe)
}

func (h *ShoeHandler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("DELETE /api/shoes/:id start userID=%d shoeID=%d", userID, id)

	if err := h.shoeService.Delete(userID, id); err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("DELETE /api/shoes/:id success userID=%d shoeID=%d", userID, id)
	response.OK(c, gin.H{"deleted": true})
}

func (h *ShoeHandler) SetPrimary(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/shoes/:id/set-primary start userID=%d shoeID=%d", userID, id)

	shoe, err := h.shoeService.SetPrimary(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/shoes/:id/set-primary success userID=%d shoeID=%d", userID, shoe.ID)
	response.OK(c, shoe)
}

func (h *ShoeHandler) Retire(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	logger.Debug("POST /api/shoes/:id/retire start userID=%d shoeID=%d", userID, id)

	shoe, err := h.shoeService.Retire(userID, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	logger.Debug("POST /api/shoes/:id/retire success userID=%d shoeID=%d", userID, shoe.ID)
	response.OK(c, shoe)
}
