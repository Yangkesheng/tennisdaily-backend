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

func (s RacketStatus) Label() string {
	switch s {
	case RacketStatusPrimary:
		return "主力拍"
	case RacketStatusActive:
		return "在用"
	case RacketStatusRetired:
		return "已退役"
	default:
		return "未知"
	}
}

type Racket struct {
	ID                         int64                    `json:"id" gorm:"primaryKey"`
	UserID                     int64                    `json:"userId" gorm:"not null;index"`
	LibraryID                  int64                    `json:"libraryId" gorm:"column:library_id;not null;default:0;index"`
	Name                       string                   `json:"name" gorm:"size:100;not null"`
	Brand                      string                   `json:"brand" gorm:"size:50"`
	Model                      string                   `json:"model" gorm:"size:100"`
	Status                     RacketStatus             `json:"status" gorm:"not null;default:2"`
	PurchaseDate               *time.Time               `json:"purchaseDate" gorm:"type:date"`
	PurchasePrice              *float64                 `json:"purchasePrice" gorm:"type:decimal(10,2)"`
	ReleaseYear                int                      `json:"releaseYear" gorm:"-"`
	Weight                     int                      `json:"weight" gorm:"-"`
	HeadSize                   int                      `json:"headSize" gorm:"-"`
	StringPattern              string                   `json:"stringPattern" gorm:"-"`
	FileID                     string                   `json:"fileId" gorm:"-"`
	StringName                 string                   `json:"stringName" gorm:"-"`
	StoreName                  string                   `json:"storeName" gorm:"-"`
	VerticalTension            *float64                 `json:"verticalTension" gorm:"-"`
	HorizontalTension          *float64                 `json:"horizontalTension" gorm:"-"`
	LastStringDate             string                   `json:"lastStringDate" gorm:"-"`
	LastStringCost             *float64                 `json:"lastStringCost" gorm:"-"`
	LatestStringingRecord      *StringingRecordResponse `json:"latestStringingRecord" gorm:"-"`
	UsageCount                 int                      `json:"usageCount" gorm:"-"`
	UsageMinutes               int                      `json:"usageMinutes" gorm:"-"`
	UsageHours                 int                      `json:"usageHours" gorm:"-"`
	AfterStringingUsageCount   int                      `json:"afterStringingUsageCount" gorm:"-"`
	AfterStringingUsageMinutes int                      `json:"afterStringingUsageMinutes" gorm:"-"`
	AfterStringingUsageHours   int                      `json:"afterStringingUsageHours" gorm:"-"`
	TotalMinutes               int                      `json:"totalMinutes" gorm:"-"`
	TotalHours                 int                      `json:"totalHours" gorm:"-"`
	CreatedAt                  time.Time                `json:"createdAt"`
	UpdatedAt                  time.Time                `json:"updatedAt"`
	DeletedAt                  gorm.DeletedAt           `json:"-" gorm:"index"`
}

func (Racket) TableName() string {
	return "racket"
}

type RacketStringingRecord struct {
	ID                int64          `json:"id" gorm:"primaryKey"`
	UserID            int64          `json:"userId" gorm:"not null;index"`
	RacketID          int64          `json:"racketId" gorm:"not null;index"`
	StringName        string         `json:"stringName" gorm:"size:100;not null"`
	StoreName         string         `json:"storeName" gorm:"column:store_name;size:100;not null;default:''"`
	VerticalTension   *float64       `json:"verticalTension" gorm:"column:vertical_tension;type:decimal(4,1)"`
	HorizontalTension *float64       `json:"horizontalTension" gorm:"column:horizontal_tension;type:decimal(4,1)"`
	Cost              float64        `json:"cost" gorm:"type:decimal(10,2);not null"`
	StringDate        time.Time      `json:"stringDate" gorm:"type:datetime;not null"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

func (RacketStringingRecord) TableName() string {
	return "racket_stringing_record"
}

type RacketUsageStats struct {
	Count   int
	Minutes int
	Hours   int
}

type RacketResponse struct {
	ID                         int64                    `json:"id"`
	LibraryID                  int64                    `json:"libraryId"`
	Name                       string                   `json:"name"`
	Brand                      string                   `json:"brand"`
	Model                      string                   `json:"model"`
	Status                     RacketStatus             `json:"status"`
	PurchaseDate               string                   `json:"purchaseDate"`
	PurchasePrice              *float64                 `json:"purchasePrice"`
	ReleaseYear                int                      `json:"releaseYear"`
	Weight                     int                      `json:"weight"`
	HeadSize                   int                      `json:"headSize"`
	StringPattern              string                   `json:"stringPattern"`
	FileID                     string                   `json:"fileId"`
	StringName                 string                   `json:"stringName"`
	StoreName                  string                   `json:"storeName"`
	VerticalTension            *float64                 `json:"verticalTension"`
	HorizontalTension          *float64                 `json:"horizontalTension"`
	LastStringDate             string                   `json:"lastStringDate"`
	LastStringCost             *float64                 `json:"lastStringCost"`
	LatestStringingRecord      *StringingRecordResponse `json:"latestStringingRecord"`
	UsageCount                 int                      `json:"usageCount"`
	UsageMinutes               int                      `json:"usageMinutes"`
	UsageHours                 int                      `json:"usageHours"`
	AfterStringingUsageCount   int                      `json:"afterStringingUsageCount"`
	AfterStringingUsageMinutes int                      `json:"afterStringingUsageMinutes"`
	AfterStringingUsageHours   int                      `json:"afterStringingUsageHours"`
	TotalMinutes               int                      `json:"totalMinutes"`
	TotalHours                 int                      `json:"totalHours"`
	CreatedAt                  time.Time                `json:"createdAt"`
	UpdatedAt                  time.Time                `json:"updatedAt"`
}

type StringingRecordResponse struct {
	ID         int64  `json:"id"`
	RacketID   int64  `json:"racketId"`
	StringName string `json:"stringName"`
	StoreName  string `json:"storeName"`

	VerticalTension   *float64  `json:"verticalTension"`
	HorizontalTension *float64  `json:"horizontalTension"`
	Cost              float64   `json:"cost"`
	StringDate        string    `json:"stringDate"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
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
	PurchaseDate  string   `json:"purchaseDate"`
	PurchasePrice *float64 `json:"purchasePrice"`
}

type UpdateRacketRequest struct {
	LibraryID     int64        `json:"libraryId"`
	Name          string       `json:"name" binding:"required"`
	Brand         string       `json:"brand"`
	Model         string       `json:"model"`
	Status        RacketStatus `json:"status"`
	PurchaseDate  string       `json:"purchaseDate"`
	PurchasePrice *float64     `json:"purchasePrice"`
}

type CreateStringingRecordRequest struct {
	StringName string `json:"stringName" binding:"required"`
	StoreName  string `json:"storeName"`

	VerticalTension   *float64 `json:"verticalTension"`
	HorizontalTension *float64 `json:"horizontalTension"`
	Cost              float64  `json:"cost"`
	StringDate        string   `json:"stringDate" binding:"required"`
}

type UpdateStringingRecordRequest struct {
	StringName string `json:"stringName" binding:"required"`
	StoreName  string `json:"storeName"`

	VerticalTension   *float64 `json:"verticalTension"`
	HorizontalTension *float64 `json:"horizontalTension"`
	Cost              float64  `json:"cost"`
	StringDate        string   `json:"stringDate" binding:"required"`
}

func NewRacketResponse(racket Racket) RacketResponse {
	purchaseDate := ""
	if racket.PurchaseDate != nil {
		purchaseDate = racket.PurchaseDate.Format("2006-01-02")
	}

	return RacketResponse{
		ID:                         racket.ID,
		LibraryID:                  racket.LibraryID,
		Name:                       racket.Name,
		Brand:                      racket.Brand,
		Model:                      racket.Model,
		Status:                     racket.Status,
		PurchaseDate:               purchaseDate,
		PurchasePrice:              racket.PurchasePrice,
		ReleaseYear:                racket.ReleaseYear,
		Weight:                     racket.Weight,
		HeadSize:                   racket.HeadSize,
		StringPattern:              racket.StringPattern,
		FileID:                     racket.FileID,
		StringName:                 racket.StringName,
		StoreName:                  racket.StoreName,
		VerticalTension:            racket.VerticalTension,
		HorizontalTension:          racket.HorizontalTension,
		LastStringDate:             racket.LastStringDate,
		LastStringCost:             racket.LastStringCost,
		LatestStringingRecord:      racket.LatestStringingRecord,
		UsageCount:                 racket.UsageCount,
		UsageMinutes:               racket.UsageMinutes,
		UsageHours:                 racket.UsageHours,
		AfterStringingUsageCount:   racket.AfterStringingUsageCount,
		AfterStringingUsageMinutes: racket.AfterStringingUsageMinutes,
		AfterStringingUsageHours:   racket.AfterStringingUsageHours,
		TotalMinutes:               racket.TotalMinutes,
		TotalHours:                 racket.TotalHours,
		CreatedAt:                  racket.CreatedAt,
		UpdatedAt:                  racket.UpdatedAt,
	}
}

func NewStringingRecordResponse(record RacketStringingRecord) StringingRecordResponse {
	return StringingRecordResponse{
		ID:                record.ID,
		RacketID:          record.RacketID,
		StringName:        record.StringName,
		StoreName:         record.StoreName,
		VerticalTension:   record.VerticalTension,
		HorizontalTension: record.HorizontalTension,
		Cost:              record.Cost,
		StringDate:        record.StringDate.Format("2006-01-02 15:04"),
		CreatedAt:         record.CreatedAt,
		UpdatedAt:         record.UpdatedAt,
	}
}
