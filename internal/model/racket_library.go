package model

import "time"

type RacketLibrary struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Brand       string    `json:"brand" gorm:"size:50;not null"`
	Model       string    `json:"model" gorm:"size:100;not null"`
	ReleaseYear int       `json:"releaseYear" gorm:"column:release_year;not null"`
	Weight      int       `json:"weight"`
	HeadSize    int       `json:"headSize" gorm:"column:head_size"`
	ImageURL    string    `json:"imageUrl" gorm:"column:image_url;size:500"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (RacketLibrary) TableName() string {
	return "racket_library"
}

type RacketLibraryItemResponse struct {
	ID          int64  `json:"id"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	ReleaseYear int    `json:"releaseYear"`
	Weight      int    `json:"weight"`
	HeadSize    int    `json:"headSize"`
	ImageURL    string `json:"imageUrl"`
}

type RacketLibraryBrandGroupResponse struct {
	Brand string                      `json:"brand"`
	Items []RacketLibraryItemResponse `json:"items"`
}

func NewRacketLibraryItemResponse(item RacketLibrary) RacketLibraryItemResponse {
	return RacketLibraryItemResponse{
		ID:          item.ID,
		Brand:       item.Brand,
		Model:       item.Model,
		ReleaseYear: item.ReleaseYear,
		Weight:      item.Weight,
		HeadSize:    item.HeadSize,
		ImageURL:    item.ImageURL,
	}
}
