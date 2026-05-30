package model

import (
	"time"

	"gorm.io/gorm"
)

type RacketStatus int16

const (
	RacketStatusPrimary RacketStatus = 1
	RacketStatusActive  RacketStatus = 2
	RacketStatusRetired RacketStatus = 3
)

type Racket struct {
	ID             int64          `json:"id" gorm:"primaryKey"`
	UserID         int64          `json:"userId" gorm:"not null;index"`
	LibraryID      int64          `json:"libraryId" gorm:"column:library_id;not null;default:0;index"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Brand          string         `json:"brand" gorm:"size:50"`
	Model          string         `json:"model" gorm:"size:100"`
	Status         RacketStatus   `json:"status" gorm:"not null;default:2"`
	ImageURL       string         `json:"imageUrl" gorm:"column:image_url;size:500"`
	PurchaseDate   *time.Time     `json:"purchaseDate" gorm:"type:date"`
	PurchasePrice  *float64       `json:"purchasePrice" gorm:"type:decimal(10,2)"`
	StringName     string         `json:"stringName" gorm:"-"`
	Tension        *float64       `json:"tension" gorm:"-"`
	LastStringDate string         `json:"lastStringDate" gorm:"-"`
	LastStringCost *float64       `json:"lastStringCost" gorm:"-"`
	TotalMinutes   int            `json:"totalMinutes" gorm:"-"`
	TotalHours     int            `json:"totalHours" gorm:"-"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Racket) TableName() string {
	return "racket"
}

type RacketStringingRecord struct {
	ID         int64          `json:"id" gorm:"primaryKey"`
	UserID     int64          `json:"userId" gorm:"not null;index"`
	RacketID   int64          `json:"racketId" gorm:"not null;index"`
	StringName string         `json:"stringName" gorm:"size:100;not null"`
	Tension    *float64       `json:"tension" gorm:"type:decimal(4,1)"`
	Cost       float64        `json:"cost" gorm:"type:decimal(10,2);not null"`
	StringDate time.Time      `json:"stringDate" gorm:"type:date;not null"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RacketStringingRecord) TableName() string {
	return "racket_stringing_record"
}

type RacketResponse struct {
	ID             int64        `json:"id"`
	LibraryID      int64        `json:"libraryId"`
	Name           string       `json:"name"`
	Brand          string       `json:"brand"`
	Model          string       `json:"model"`
	Status         RacketStatus `json:"status"`
	ImageURL       string       `json:"imageUrl"`
	PurchaseDate   string       `json:"purchaseDate"`
	PurchasePrice  *float64     `json:"purchasePrice"`
	StringName     string       `json:"stringName"`
	Tension        *float64     `json:"tension"`
	LastStringDate string       `json:"lastStringDate"`
	LastStringCost *float64     `json:"lastStringCost"`
	TotalMinutes   int          `json:"totalMinutes"`
	TotalHours     int          `json:"totalHours"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

type StringingRecordResponse struct {
	ID         int64     `json:"id"`
	RacketID   int64     `json:"racketId"`
	StringName string    `json:"stringName"`
	Tension    *float64  `json:"tension"`
	Cost       float64   `json:"cost"`
	StringDate string    `json:"stringDate"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type RacketDetailResponse struct {
	Racket           RacketResponse            `json:"racket"`
	StringingRecords []StringingRecordResponse `json:"stringingRecords"`
}

type RacketStatsResponse struct {
	RacketCount       int     `json:"racketCount"`
	RacketCost        float64 `json:"racketCost"`
	StringingCost     float64 `json:"stringingCost"`
	TotalCost         float64 `json:"totalCost"`
	RacketCostText    string  `json:"racketCostText"`
	StringingCostText string  `json:"stringingCostText"`
	TotalCostText     string  `json:"totalCostText"`
}

type CreateRacketRequest struct {
	LibraryID     int64    `json:"libraryId"`
	Name          string   `json:"name" binding:"required"`
	Brand         string   `json:"brand"`
	Model         string   `json:"model"`
	ImageURL      string   `json:"imageUrl"`
	PurchaseDate  string   `json:"purchaseDate"`
	PurchasePrice *float64 `json:"purchasePrice"`
}

type UpdateRacketRequest struct {
	LibraryID     int64        `json:"libraryId"`
	Name          string       `json:"name" binding:"required"`
	Brand         string       `json:"brand"`
	Model         string       `json:"model"`
	Status        RacketStatus `json:"status"`
	ImageURL      string       `json:"imageUrl"`
	PurchaseDate  string       `json:"purchaseDate"`
	PurchasePrice *float64     `json:"purchasePrice"`
}

type CreateStringingRecordRequest struct {
	StringName string   `json:"stringName" binding:"required"`
	Tension    *float64 `json:"tension"`
	Cost       float64  `json:"cost"`
	StringDate string   `json:"stringDate" binding:"required"`
}

func NewRacketResponse(racket Racket) RacketResponse {
	purchaseDate := ""
	if racket.PurchaseDate != nil {
		purchaseDate = racket.PurchaseDate.Format("2006-01-02")
	}

	return RacketResponse{
		ID:             racket.ID,
		LibraryID:      racket.LibraryID,
		Name:           racket.Name,
		Brand:          racket.Brand,
		Model:          racket.Model,
		Status:         racket.Status,
		ImageURL:       racket.ImageURL,
		PurchaseDate:   purchaseDate,
		PurchasePrice:  racket.PurchasePrice,
		StringName:     racket.StringName,
		Tension:        racket.Tension,
		LastStringDate: racket.LastStringDate,
		LastStringCost: racket.LastStringCost,
		TotalMinutes:   racket.TotalMinutes,
		TotalHours:     racket.TotalHours,
		CreatedAt:      racket.CreatedAt,
		UpdatedAt:      racket.UpdatedAt,
	}
}

func NewStringingRecordResponse(record RacketStringingRecord) StringingRecordResponse {
	return StringingRecordResponse{
		ID:         record.ID,
		RacketID:   record.RacketID,
		StringName: record.StringName,
		Tension:    record.Tension,
		Cost:       record.Cost,
		StringDate: record.StringDate.Format("2006-01-02"),
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
	}
}
