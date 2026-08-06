package model

import (
	"time"

	"gorm.io/gorm"
)

type ShoeStatus int16

const (
	ShoeStatusPrimary ShoeStatus = 1
	ShoeStatusActive  ShoeStatus = 2
	ShoeStatusRetired ShoeStatus = 3
)

func (s ShoeStatus) Label() string {
	switch s {
	case ShoeStatusPrimary:
		return "主力鞋"
	case ShoeStatusActive:
		return "在用"
	case ShoeStatusRetired:
		return "已退役"
	default:
		return "未知"
	}
}

type Shoe struct {
	ID            int64             `json:"id" gorm:"primaryKey"`
	UserID        int64             `json:"userId" gorm:"not null;index"`
	LibraryID     int64             `json:"libraryId" gorm:"column:library_id;not null;default:0;index"`
	Name          string            `json:"name" gorm:"size:100;not null"`
	Brand         string            `json:"brand" gorm:"size:50"`
	Model         string            `json:"model" gorm:"size:100"`
	Status        ShoeStatus        `json:"status" gorm:"not null;default:2"`
	Size          string            `json:"size" gorm:"size:20"`
	Colorway      string            `json:"colorway" gorm:"size:128"`
	PurchaseDate  *time.Time        `json:"purchaseDate" gorm:"type:date"`
	PurchasePrice *float64          `json:"purchasePrice" gorm:"type:decimal(10,2)"`
	Gender        int               `json:"gender" gorm:"-"`
	ReleaseYear   int               `json:"releaseYear" gorm:"-"`
	FileID        string            `json:"fileId" gorm:"-"`
	UsageCount    int               `json:"usageCount" gorm:"-"`
	UsageMinutes  int               `json:"usageMinutes" gorm:"-"`
	UsageHours    int               `json:"usageHours" gorm:"-"`
	TotalMinutes  int               `json:"totalMinutes" gorm:"-"`
	TotalHours    int               `json:"totalHours" gorm:"-"`
	Wear          *ShoeWearResponse `json:"wear" gorm:"-"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
}

func (Shoe) TableName() string {
	return "my_shoes"
}

type ShoeBrand struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	Slug      string    `json:"-" gorm:"size:128;not null"`
	FileID    string    `json:"fileId" gorm:"column:file_id;size:255"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ShoeBrand) TableName() string {
	return "shoe_brands"
}

type ShoeSeries struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	BrandID   int64     `json:"brandId" gorm:"column:brand_id;not null"`
	Gender    int       `json:"gender" gorm:"not null;default:0"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ShoeSeries) TableName() string {
	return "shoe_series"
}

type ShoeLibrary struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	BrandID       int64     `json:"brandId" gorm:"column:brand_id;not null;default:0"`
	Brand         string    `json:"brand" gorm:"size:50;not null"`
	SeriesID      int64     `json:"seriesId" gorm:"column:series_id;not null;default:0"`
	Series        string    `json:"series" gorm:"size:100;not null;default:''"`
	Model         string    `json:"model" gorm:"size:100;not null"`
	Gender        int       `json:"gender" gorm:"not null;default:0"`
	Colorway      string    `json:"colorway" gorm:"size:128"`
	ProductCode   string    `json:"-" gorm:"column:product_code;size:64"`
	ReleaseYear   int       `json:"releaseYear" gorm:"column:release_year;not null;default:0"`
	Weight        string    `json:"weight" gorm:"size:128"`
	Width         string    `json:"width" gorm:"size:128"`
	Surface       string    `json:"surface" gorm:"size:255"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
	ColorwayCount int       `json:"colorwayCount" gorm:"column:colorway_count;not null;default:0"`
	FileID        string    `json:"fileId" gorm:"column:file_id;size:255"`
	ImageURL      string    `json:"imageUrl" gorm:"column:image_url;size:500"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (ShoeLibrary) TableName() string {
	return "shoe_library"
}

// ShoeLibraryImageUpload 记录球鞋库图片已上传对象存储的日志，用于迁移幂等与审计回滚。
type ShoeLibraryImageUpload struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	ShoeLibraryID int64     `json:"shoeLibraryId" gorm:"column:shoe_library_id;not null;uniqueIndex"`
	FileID        string    `json:"fileId" gorm:"column:file_id;size:255;not null"`
	ObjectKey     string    `json:"objectKey" gorm:"column:object_key;size:255;not null"`
	SourceURL     string    `json:"sourceUrl" gorm:"column:source_url;size:500;not null;default:''"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (ShoeLibraryImageUpload) TableName() string {
	return "shoe_library_image_upload"
}

type ShoeSeriesQuery struct {
	BrandID int64 `form:"brandId" binding:"omitempty,min=1"`
	Gender  *int  `form:"gender" binding:"omitempty,min=0,max=3"`
}

type ShoeLibraryQuery struct {
	BrandID  int64 `form:"brandId" binding:"omitempty,min=1"`
	SeriesID int64 `form:"seriesId" binding:"omitempty,min=1"`
	Gender   *int  `form:"gender" binding:"omitempty,min=0,max=3"`
}

type ShoeBrandResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	FileID string `json:"fileId"`
}

type ShoeSeriesResponse struct {
	ID      int64  `json:"id"`
	BrandID int64  `json:"brandId"`
	Gender  int    `json:"gender"`
	Name    string `json:"name"`
}

type ShoeLibraryItemResponse struct {
	ID            int64   `json:"id"`
	BrandID       int64   `json:"brandId"`
	Brand         string  `json:"brand"`
	SeriesID      int64   `json:"seriesId"`
	Series        string  `json:"series"`
	Model         string  `json:"model"`
	Gender        int     `json:"gender"`
	ReleaseYear   int     `json:"releaseYear"`
	Colorway      string  `json:"colorway"`
	Weight        string  `json:"weight"`
	Width         string  `json:"width"`
	Surface       string  `json:"surface"`
	Price         float64 `json:"price"`
	ColorwayCount int     `json:"colorwayCount"`
	FileID        string  `json:"fileId"`
	ImageURL      string  `json:"imageUrl"`
}

type ShoeLibraryBrandGroupResponse struct {
	BrandID int64                     `json:"brandId"`
	Brand   string                    `json:"brand"`
	Items   []ShoeLibraryItemResponse `json:"items"`
}

type ShoeLibraryBrandStats struct {
	BrandID int64  `gorm:"column:brand_id"`
	Brand   string `gorm:"column:brand"`
	Count   int    `gorm:"column:count"`
}

type ShoeLibrarySeriesStats struct {
	BrandID  int64  `gorm:"column:brand_id"`
	Brand    string `gorm:"column:brand"`
	SeriesID int64  `gorm:"column:series_id"`
	Series   string `gorm:"column:series"`
	Gender   int    `gorm:"column:gender"`
	Count    int    `gorm:"column:count"`
}

type ShoeLibraryBrandStatsResponse struct {
	BrandID int64                                       `json:"brandId"`
	Brand   string                                      `json:"brand"`
	Count   int                                         `json:"count"`
	Series  map[string][]ShoeLibrarySeriesStatsResponse `json:"series"`
}

type ShoeLibrarySeriesStatsResponse struct {
	SeriesID int64  `json:"seriesId"`
	Series   string `json:"series"`
	Count    int    `json:"count"`
}

type ShoeUsageStats struct {
	Count   int
	Minutes int
	Hours   int
}

type ShoeWearResponse struct {
	State          string  `json:"state"`
	Display        string  `json:"display"`
	Score          float64 `json:"score"`
	RemainingHours int     `json:"remainingHours"`
}

type ShoeResponse struct {
	ID            int64             `json:"id"`
	LibraryID     int64             `json:"libraryId"`
	Name          string            `json:"name"`
	Brand         string            `json:"brand"`
	Model         string            `json:"model"`
	Status        ShoeStatus        `json:"status"`
	Gender        int               `json:"gender"`
	Size          string            `json:"size"`
	Colorway      string            `json:"colorway"`
	PurchaseDate  string            `json:"purchaseDate"`
	PurchasePrice *float64          `json:"purchasePrice"`
	ReleaseYear   int               `json:"releaseYear"`
	FileID        string            `json:"fileId"`
	UsageCount    int               `json:"usageCount"`
	UsageMinutes  int               `json:"usageMinutes"`
	UsageHours    int               `json:"usageHours"`
	TotalMinutes  int               `json:"totalMinutes"`
	TotalHours    int               `json:"totalHours"`
	Wear          *ShoeWearResponse `json:"wear"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type ShoeStatsResponse struct {
	ShoeCount     int     `json:"shoeCount"`
	ShoeCost      float64 `json:"shoeCost"`
	TotalCost     float64 `json:"totalCost"`
	ShoeCostText  string  `json:"shoeCostText"`
	TotalCostText string  `json:"totalCostText"`
}

type CreateShoeRequest struct {
	LibraryID     int64      `json:"libraryId"`
	Name          string     `json:"name" binding:"required"`
	Brand         string     `json:"brand"`
	Model         string     `json:"model"`
	Status        ShoeStatus `json:"status"`
	Size          string     `json:"size"`
	Colorway      string     `json:"colorway"`
	PurchaseDate  string     `json:"purchaseDate"`
	PurchasePrice *float64   `json:"purchasePrice"`
}

type UpdateShoeRequest struct {
	LibraryID     int64      `json:"libraryId"`
	Name          string     `json:"name" binding:"required"`
	Brand         string     `json:"brand"`
	Model         string     `json:"model"`
	Status        ShoeStatus `json:"status"`
	Size          string     `json:"size"`
	Colorway      string     `json:"colorway"`
	PurchaseDate  string     `json:"purchaseDate"`
	PurchasePrice *float64   `json:"purchasePrice"`
}

// ---- 管理员维护球鞋库请求/响应 ----

type CreateShoeBrandRequest struct {
	Name   string `json:"name" binding:"required"`
	Slug   string `json:"slug"`
	FileID string `json:"fileId"`
}

type CreateShoeSeriesRequest struct {
	BrandID int64  `json:"brandId" binding:"required,min=1"`
	Gender  int    `json:"gender" binding:"omitempty,min=0,max=3"`
	Name    string `json:"name" binding:"required"`
}

type CreateShoeLibraryRequest struct {
	BrandID       int64   `json:"brandId" binding:"required,min=1"`
	SeriesID      int64   `json:"seriesId" binding:"required,min=1"`
	Model         string  `json:"model" binding:"required"`
	Gender        int     `json:"gender" binding:"omitempty,min=0,max=3"`
	Colorway      string  `json:"colorway"`
	ReleaseYear   int     `json:"releaseYear"`
	Weight        string  `json:"weight"`
	Width         string  `json:"width"`
	Surface       string  `json:"surface"`
	Price         float64 `json:"price"`
	ColorwayCount int     `json:"colorwayCount"`
	FileID        string  `json:"fileId"`
	ImageURL      string  `json:"imageUrl"`
}

type CreateShoeLibraryResponse struct {
	Created int                       `json:"created"`
	Items   []ShoeLibraryItemResponse `json:"items"`
}

func NewShoeBrandResponse(item ShoeBrand) ShoeBrandResponse {
	return ShoeBrandResponse{
		ID:     item.ID,
		Name:   item.Name,
		FileID: item.FileID,
	}
}

func NewShoeSeriesResponse(item ShoeSeries) ShoeSeriesResponse {
	return ShoeSeriesResponse{
		ID:      item.ID,
		BrandID: item.BrandID,
		Gender:  item.Gender,
		Name:    item.Name,
	}
}

func NewShoeLibraryItemResponse(item ShoeLibrary) ShoeLibraryItemResponse {
	return ShoeLibraryItemResponse{
		ID:            item.ID,
		BrandID:       item.BrandID,
		Brand:         item.Brand,
		SeriesID:      item.SeriesID,
		Series:        item.Series,
		Model:         item.Model,
		Gender:        item.Gender,
		ReleaseYear:   item.ReleaseYear,
		Colorway:      item.Colorway,
		Weight:        item.Weight,
		Width:         item.Width,
		Surface:       item.Surface,
		Price:         item.Price,
		ColorwayCount: item.ColorwayCount,
		FileID:        item.FileID,
		ImageURL:      item.ImageURL,
	}
}

func NewShoeResponse(shoe Shoe) ShoeResponse {
	purchaseDate := ""
	if shoe.PurchaseDate != nil {
		purchaseDate = shoe.PurchaseDate.Format("2006-01-02")
	}

	return ShoeResponse{
		ID:            shoe.ID,
		LibraryID:     shoe.LibraryID,
		Name:          shoe.Name,
		Brand:         shoe.Brand,
		Model:         shoe.Model,
		Status:        shoe.Status,
		Gender:        shoe.Gender,
		Size:          shoe.Size,
		Colorway:      shoe.Colorway,
		PurchaseDate:  purchaseDate,
		PurchasePrice: shoe.PurchasePrice,
		ReleaseYear:   shoe.ReleaseYear,
		FileID:        shoe.FileID,
		UsageCount:    shoe.UsageCount,
		UsageMinutes:  shoe.UsageMinutes,
		UsageHours:    shoe.UsageHours,
		TotalMinutes:  shoe.TotalMinutes,
		TotalHours:    shoe.TotalHours,
		Wear:          shoe.Wear,
		CreatedAt:     shoe.CreatedAt,
		UpdatedAt:     shoe.UpdatedAt,
	}
}
