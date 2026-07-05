package model

import (
	"time"

	"gorm.io/gorm"
)

type TennisSession struct {
	ID              int64              `json:"id" gorm:"primaryKey"`
	UserID          int64              `json:"userId" gorm:"not null;index"`
	Date            time.Time          `json:"date" gorm:"type:datetime;not null"`
	DurationMinutes int                `json:"durationMinutes" gorm:"not null;default:120"`
	Rating          int16              `json:"rating" gorm:"not null;default:3"`
	Type            SessionType        `json:"type" gorm:"not null;default:1"`
	Category        SessionCategory    `json:"category" gorm:"not null;default:1"`
	SubCategory     SessionSubCategory `json:"subCategory" gorm:"column:sub_category;not null;default:2"`
	MatchRank       MatchRank          `json:"matchRank" gorm:"not null;default:0"`
	CourtName       string             `json:"courtName" gorm:"size:128;not null;default:''"`
	Partner         string             `json:"partner" gorm:"size:128;not null;default:''"`
	Cost            float64            `json:"cost" gorm:"type:decimal(10,2);not null;default:0"`
	RacketID        int64              `json:"racketId" gorm:"not null;default:0;index"`
	RacketName      string             `json:"racketName" gorm:"size:128;not null;default:''"`
	ShoeName        string             `json:"shoeName" gorm:"size:128;not null;default:''"`
	Note            string             `json:"note" gorm:"not null;default:''"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt     `json:"-" gorm:"index"`
}

func (TennisSession) TableName() string {
	return "tennis_sessions"
}

type SessionResponse struct {
	ID              int64              `json:"id"`
	Date            string             `json:"date"`
	DurationMinutes int                `json:"durationMinutes"`
	Rating          int16              `json:"rating"`
	Type            SessionType        `json:"type"`
	TypeLabel       string             `json:"typeLabel"`
	Category        SessionCategory    `json:"category"`
	CategoryText    string             `json:"categoryText"`
	SubCategory     SessionSubCategory `json:"subCategory"`
	SubCategoryText string             `json:"subCategoryText"`
	TypeText        string             `json:"typeText"`
	MatchRank       MatchRank          `json:"matchRank"`
	MatchRankLabel  string             `json:"matchRankLabel"`
	CourtName       string             `json:"courtName"`
	Partner         string             `json:"partner"`
	Cost            float64            `json:"cost"`
	RacketID        int64              `json:"racketId"`
	RacketName      string             `json:"racketName"`
	ShoeName        string             `json:"shoeName"`
	Note            string             `json:"note"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

type SessionCalendarDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type SessionCalendarSummaryResponse struct {
	SessionCount   int64   `json:"sessionCount"`
	ActiveDayCount int64   `json:"activeDayCount"`
	TotalMinutes   int64   `json:"totalMinutes"`
	AverageMinutes float64 `json:"averageMinutes"`
	AverageRating  float64 `json:"averageRating"`
	SessionCost    float64 `json:"sessionCost"`
	RacketCost     float64 `json:"racketCost"`
	StringingCost  float64 `json:"stringingCost"`
	TotalCost      float64 `json:"totalCost"`
	TrainingCount  int64   `json:"trainingCount"`
	SinglesCount   int64   `json:"singlesCount"`
	DoublesCount   int64   `json:"doublesCount"`
	MatchCount     int64   `json:"matchCount"`
}

type CalendarWeeklySessionChartItemResponse struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type CalendarRatingTrendChartItemResponse struct {
	Date   string `json:"date"`
	Rating int16  `json:"rating"`
}

type CalendarExpenseChartItemResponse struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type SessionCalendarChartsResponse struct {
	WeeklySessions       []CalendarWeeklySessionChartItemResponse `json:"weeklySessions"`
	RatingTrend          []CalendarRatingTrendChartItemResponse   `json:"ratingTrend"`
	ExpenseBreakdown     []CalendarExpenseChartItemResponse       `json:"expenseBreakdown"`
	SessionTypeBreakdown []StatsBreakdownItemResponse             `json:"sessionTypeBreakdown"`
}

type SessionCalendarResponse struct {
	Year           int                            `json:"year"`
	Month          int                            `json:"month"`
	ActiveDayCount int                            `json:"activeDayCount"`
	Days           []SessionCalendarDay           `json:"days"`
	Summary        SessionCalendarSummaryResponse `json:"summary"`
	Charts         SessionCalendarChartsResponse  `json:"charts"`
}

type SessionCategoryTextResolver interface {
	Label(category SessionCategory) (string, bool)
	SubCategoryLabel(category SessionCategory, subCategory SessionSubCategory) (string, bool)
	TypeText(category SessionCategory, subCategory SessionSubCategory) (string, bool)
}

func NewSessionResponse(session TennisSession, resolver SessionCategoryTextResolver) SessionResponse {
	category, subCategory := session.Category, session.SubCategory
	categoryText, subCategoryText, typeText, ok := sessionCategoryTexts(category, subCategory, resolver)
	if !ok {
		category, subCategory = session.Type.ToCategoryPair()
		categoryText, subCategoryText, typeText, _ = sessionCategoryTexts(category, subCategory, resolver)
	}

	return SessionResponse{
		ID:              session.ID,
		Date:            session.Date.Format("2006-01-02 15:04"),
		DurationMinutes: session.DurationMinutes,
		Rating:          session.Rating,
		Type:            session.Type,
		TypeLabel:       session.Type.Label(),
		Category:        category,
		CategoryText:    categoryText,
		SubCategory:     subCategory,
		SubCategoryText: subCategoryText,
		TypeText:        typeText,
		MatchRank:       session.MatchRank,
		MatchRankLabel:  session.MatchRank.Label(),
		CourtName:       session.CourtName,
		Partner:         session.Partner,
		Cost:            session.Cost,
		RacketID:        session.RacketID,
		RacketName:      session.RacketName,
		ShoeName:        session.ShoeName,
		Note:            session.Note,
		CreatedAt:       session.CreatedAt,
		UpdatedAt:       session.UpdatedAt,
	}
}

func sessionCategoryTexts(category SessionCategory, subCategory SessionSubCategory, resolver SessionCategoryTextResolver) (string, string, string, bool) {
	if resolver != nil {
		categoryText, categoryOK := resolver.Label(category)
		subCategoryText, subCategoryOK := resolver.SubCategoryLabel(category, subCategory)
		typeText, typeTextOK := resolver.TypeText(category, subCategory)
		if categoryOK && subCategoryOK && typeTextOK {
			return categoryText, subCategoryText, typeText, true
		}
	}

	if !category.IsValid() || !category.IsValidSubCategory(subCategory) {
		return "", "", "", false
	}
	return category.Label(), subCategory.Label(category), category.TypeText(subCategory), true
}

type SessionListQuery struct {
	Date     string `form:"date"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type SessionListPageResponse struct {
	List       []SessionResponse `json:"list"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
	HasMore    bool              `json:"hasMore"`
}

type CreateSessionRequest struct {
	Date            string             `json:"date" binding:"required"`
	DurationMinutes int                `json:"durationMinutes"`
	Rating          int16              `json:"rating"`
	Type            SessionType        `json:"type"`
	Category        SessionCategory    `json:"category"`
	SubCategory     SessionSubCategory `json:"subCategory"`
	MatchRank       MatchRank          `json:"matchRank"`
	CourtName       string             `json:"courtName"`
	Partner         string             `json:"partner"`
	Cost            float64            `json:"cost"`
	RacketID        int64              `json:"racketId"`
	RacketName      string             `json:"racketName"`
	ShoeName        string             `json:"shoeName"`
	Note            string             `json:"note"`
}

type UpdateSessionRequest struct {
	Date            string             `json:"date" binding:"required"`
	DurationMinutes int                `json:"durationMinutes"`
	Rating          int16              `json:"rating"`
	Type            SessionType        `json:"type"`
	Category        SessionCategory    `json:"category"`
	SubCategory     SessionSubCategory `json:"subCategory"`
	MatchRank       MatchRank          `json:"matchRank"`
	CourtName       string             `json:"courtName"`
	Partner         string             `json:"partner"`
	Cost            float64            `json:"cost"`
	RacketID        int64              `json:"racketId"`
	RacketName      string             `json:"racketName"`
	ShoeName        string             `json:"shoeName"`
	Note            string             `json:"note"`
}
