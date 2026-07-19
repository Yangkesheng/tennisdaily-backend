package model

import "time"

type RacketBrand struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	FileID    string    `json:"fileId" gorm:"column:file_id;size:255"`
	CreatedAt time.Time `json:"createdAt"`
}

func (RacketBrand) TableName() string {
	return "racket_brands"
}

type RacketSeries struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	BrandID   int64     `json:"brandId" gorm:"column:brand_id;not null"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

func (RacketSeries) TableName() string {
	return "racket_series"
}

type RacketLibrary struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	BrandID       int64     `json:"brandId" gorm:"column:brand_id;not null;default:0"`
	Brand         string    `json:"brand" gorm:"size:50;not null"`
	SeriesID      int64     `json:"seriesId" gorm:"column:series_id;not null;default:0"`
	Series        string    `json:"series" gorm:"size:100;not null;default:''"`
	Model         string    `json:"model" gorm:"size:100;not null"`
	ReleaseYear   int       `json:"releaseYear" gorm:"column:release_year;not null"`
	Weight        int       `json:"weight"`
	HeadSize      int       `json:"headSize" gorm:"column:head_size"`
	StringPattern string    `json:"stringPattern" gorm:"column:string_pattern;size:64"`
	FileID        string    `json:"fileId" gorm:"column:file_id;size:255"`
	ImageURL      string    `json:"imageUrl" gorm:"column:image_url;size:500"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (RacketLibrary) TableName() string {
	return "racket_library"
}

type RacketLibraryQuery struct {
	BrandID  int64 `form:"brandId" binding:"omitempty,min=1"`
	SeriesID int64 `form:"seriesId" binding:"omitempty,min=1"`
}

type RacketSeriesQuery struct {
	BrandID int64 `form:"brandId" binding:"omitempty,min=1"`
}

type RacketBrandResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	FileID string `json:"fileId"`
}

type RacketSeriesResponse struct {
	ID      int64  `json:"id"`
	BrandID int64  `json:"brandId"`
	Name    string `json:"name"`
}

func NewRacketBrandResponse(item RacketBrand) RacketBrandResponse {
	return RacketBrandResponse{
		ID:     item.ID,
		Name:   item.Name,
		FileID: item.FileID,
	}
}

func NewRacketSeriesResponse(item RacketSeries) RacketSeriesResponse {
	return RacketSeriesResponse{
		ID:      item.ID,
		BrandID: item.BrandID,
		Name:    item.Name,
	}
}

type RacketLibraryItemResponse struct {
	ID            int64  `json:"id"`
	BrandID       int64  `json:"brandId"`
	Brand         string `json:"brand"`
	SeriesID      int64  `json:"seriesId"`
	Series        string `json:"series"`
	Model         string `json:"model"`
	ReleaseYear   int    `json:"releaseYear"`
	Weight        int    `json:"weight"`
	HeadSize      int    `json:"headSize"`
	StringPattern string `json:"stringPattern"`
	FileID        string `json:"fileId"`
	ImageURL      string `json:"imageUrl"`
}

type RacketLibraryBrandGroupResponse struct {
	BrandID int64                       `json:"brandId"`
	Brand   string                      `json:"brand"`
	Items   []RacketLibraryItemResponse `json:"items"`
}

func NewRacketLibraryItemResponse(item RacketLibrary) RacketLibraryItemResponse {
	return RacketLibraryItemResponse{
		ID:            item.ID,
		BrandID:       item.BrandID,
		Brand:         item.Brand,
		SeriesID:      item.SeriesID,
		Series:        item.Series,
		Model:         item.Model,
		ReleaseYear:   item.ReleaseYear,
		Weight:        item.Weight,
		HeadSize:      item.HeadSize,
		StringPattern: item.StringPattern,
		FileID:        item.FileID,
		ImageURL:      item.ImageURL,
	}
}
