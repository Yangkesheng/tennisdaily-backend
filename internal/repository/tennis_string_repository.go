package repository

import (
	"gorm.io/gorm"
	"strings"
	"tennisdaily-backend/internal/model"
)

type TennisStringRepository struct{ db *gorm.DB }

func NewTennisStringRepository(db *gorm.DB) *TennisStringRepository {
	return &TennisStringRepository{db: db}
}
func stringSearchPattern(keyword string) string {
	return "%" + strings.NewReplacer("!", "!!", "%", "!%", "_", "!_").Replace(keyword) + "%"
}
func (r *TennisStringRepository) List(q model.StringLibraryQuery) (model.StringLibraryPage, error) {
	result := model.StringLibraryPage{Items: []model.StringLibrary{}, Page: q.Page, PageSize: q.PageSize}
	db := r.db.Model(&model.StringLibrary{}).Where("published = ?", true)
	if q.StringID > 0 {
		db = db.Where("id = ?", q.StringID)
	}
	if q.Brand != "" {
		db = db.Where("brand = ?", q.Brand)
	}
	if q.Keyword != "" {
		pattern := stringSearchPattern(strings.ToLower(q.Keyword))
		db = db.Where("LOWER(CONCAT(brand, ' ', model, ' ', name_zh, ' ', aliases)) LIKE ? ESCAPE '!'", pattern)
	}
	if err := db.Count(&result.Total).Error; err != nil {
		return result, err
	}
	err := db.Preload("Variants", func(db *gorm.DB) *gorm.DB { return db.Order("diameter_mm ASC") }).Order("brand ASC, model ASC, id ASC").Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&result.Items).Error
	return result, err
}
func (r *TennisStringRepository) Brands() ([]string, error) {
	items := []string{}
	err := r.db.Model(&model.StringLibrary{}).Where("published = ?", true).Distinct().Order("brand ASC").Pluck("brand", &items).Error
	return items, err
}
